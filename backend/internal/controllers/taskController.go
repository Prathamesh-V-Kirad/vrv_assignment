package controllers

import (
	"context"
	"errors"
	"backend/internal/models"
	"backend/internal/database"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)


func ParseJWT(token string) (*models.CustomClaims, error) {
	secretKey := os.Getenv("JWT_SECRET_KEY")

	claims := &models.CustomClaims{}
    parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("unexpected signing method")
        }
        return []byte(secretKey), nil
    })

    if err != nil {
        return nil, err
    }

    
    if !parsedToken.Valid {
        return nil, errors.New("invalid token")
    }

    return claims, nil
}


func getUserIDFromToken(c *fiber.Ctx) (primitive.ObjectID, error) {
	
	token := c.Cookies("jwt")
	if token == "" {
		return primitive.ObjectID{}, errors.New("Authorization token is missing")
	}

	
	
	claims, err := ParseJWT(token)
	if err != nil {
		return primitive.ObjectID{}, err
	}

	
	userID, err := primitive.ObjectIDFromHex(claims.Issuer) 
	if err != nil {
		return primitive.ObjectID{}, errors.New("Invalid user ID in token")
	}

	return userID, nil
}

func GetTasks(c *fiber.Ctx) error {
    userID, err := getUserIDFromToken(c)
    if err != nil {
        return err
    }

    
    userCollection := database.GetCollection("users")
    var user models.User
    err = userCollection.FindOne(context.Background(), bson.M{"_id": userID}).Decode(&user)
    if err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
    }

    
    roleCollection := database.GetCollection("roles")
    var role models.Role
    err = roleCollection.FindOne(context.Background(), bson.M{"_id": user.RoleID}).Decode(&role)
    if err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Role not found"})
    }

    
    permissionCollection := database.GetCollection("permissions")
    var permissions []models.Permission
    cursor, err := permissionCollection.Find(context.Background(), bson.M{"_id": bson.M{"$in": role.Permissions}})
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve permissions"})
    }
    defer cursor.Close(context.Background())

    
    if err := cursor.All(context.Background(), &permissions); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to decode permissions"})
    }

    hasPermission := false
    for _, permission := range permissions {
        if permission.Name == "view_task" { 
            hasPermission = true
            break
        }
    }

    if !hasPermission {
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "You do not have permission to view tasks"})
    }

    // Fetch the tasks from the database
    taskCollection := database.GetCollection("tasks")
    cursor, err = taskCollection.Find(context.Background(), bson.M{})
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve tasks"})
    }
    defer cursor.Close(context.Background())

    var tasks []models.Task
    if err := cursor.All(context.Background(), &tasks); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to decode tasks"})
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{"tasks": tasks})
}


func CreateTask(c *fiber.Ctx) error {
    var task models.Task
    if err := c.BodyParser(&task); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
    }

    
    token := c.Cookies("jwt")
    if token == "" {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Authorization token is missing"})
    }

    
    claims, err := ParseJWT(token)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
    }

    userID, err := primitive.ObjectIDFromHex(claims.Issuer)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid user ID in token",
		})
	}

    
    userCollection := database.GetCollection("users")
    var user models.User
    err = userCollection.FindOne(context.Background(), bson.M{"_id": userID}).Decode(&user)
    if err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
    }

    
    roleCollection := database.GetCollection("roles")
    var role models.Role
    err = roleCollection.FindOne(context.Background(), bson.M{"_id": user.RoleID}).Decode(&role)
    if err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Role not found"})
    }

    
    permissionCollection := database.GetCollection("permissions")
    var permissions []models.Permission
    cursor, err := permissionCollection.Find(context.Background(), bson.M{"_id": bson.M{"$in": role.Permissions}})
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve permissions"})
    }
    defer cursor.Close(context.Background())

    if err := cursor.All(context.Background(), &permissions); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to decode permissions"})
    }


    hasPermission := false
    for _, permission := range permissions {
        if permission.Name == "create_task" { 
            hasPermission = true
            break
        }
    }

    if !hasPermission {
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "You do not have permission to create a task"})
    }

    
    task.CreatedAt = time.Now()
    task.UpdatedAt = time.Now()
    collection := database.GetCollection("tasks")
    _, err = collection.InsertOne(context.Background(), task)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create task"})
    }

    return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Task created successfully", "task": task})
}


func UpdateTask(c *fiber.Ctx) error {
	taskID := c.Params("id")
	var taskUpdate models.Task
	if err := c.BodyParser(&taskUpdate); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	userID, err := getUserIDFromToken(c)
	if err != nil {
		return err
	}

	
	userCollection := database.GetCollection("users")
	var user models.User
	err = userCollection.FindOne(context.Background(), bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	
	roleCollection := database.GetCollection("roles")
	var role models.Role
	err = roleCollection.FindOne(context.Background(), bson.M{"_id": user.RoleID}).Decode(&role)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Role not found"})
	}

	
	permissionCollection := database.GetCollection("permissions")
	var permissions []models.Permission
	cursor, err := permissionCollection.Find(context.Background(), bson.M{"_id": bson.M{"$in": role.Permissions}})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve permissions"})
	}
	defer cursor.Close(context.Background())

	if err := cursor.All(context.Background(), &permissions); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to decode permissions"})
	}

	
	hasPermission := false
	for _, permission := range permissions {
		if permission.Name == "update_task_status" { 
			hasPermission = true
			break
		}
	}

	if !hasPermission {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "You do not have permission to update this task"})
	}

	
	taskObjectID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid task ID"})
	}

	
	update := bson.M{"$set": taskUpdate}
	collection := database.GetCollection("tasks")
	_, err = collection.UpdateOne(context.Background(), bson.M{"_id": taskObjectID}, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update task"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Task updated successfully"})
}

func DeleteTask(c *fiber.Ctx) error {
	taskID := c.Params("id")

	userID, err := getUserIDFromToken(c)
	if err != nil {
		return err
	}

	
	userCollection := database.GetCollection("users")
	var user models.User
	err = userCollection.FindOne(context.Background(), bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	
	roleCollection := database.GetCollection("roles")
	var role models.Role
	err = roleCollection.FindOne(context.Background(), bson.M{"_id": user.RoleID}).Decode(&role)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Role not found"})
	}

	
	permissionCollection := database.GetCollection("permissions")
	var permissions []models.Permission
	cursor, err := permissionCollection.Find(context.Background(), bson.M{"_id": bson.M{"$in": role.Permissions}})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve permissions"})
	}
	defer cursor.Close(context.Background())

	if err := cursor.All(context.Background(), &permissions); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to decode permissions"})
	}

	
	hasPermission := false
	for _, permission := range permissions {
		if permission.Name == "delete_task" { 
			hasPermission = true
			break
		}
	}

	if !hasPermission {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "You do not have permission to delete this task"})
	}

	
	taskObjectID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid task ID"})
	}

	collection := database.GetCollection("tasks")
	_, err = collection.DeleteOne(context.Background(), bson.M{"_id": taskObjectID})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete task"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Task deleted successfully"})
}
