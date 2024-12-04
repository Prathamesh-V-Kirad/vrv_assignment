package initialize

import (
	"log"
	"backend/internal/database"
	"context"
	"backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func InitializePermissionsAndRoles() {
	
	permissions := []models.Permission{
		{Name: "view_task", Description: "Allows viewing tasks"},
		{Name: "create_task", Description: "Allows creating tasks"},
		{Name: "update_task", Description: "Allows updating tasks"},
		{Name: "delete_task", Description: "Allows deleting tasks"},
	}

	
	roles := []models.Role{
		{Name: "admin", Permissions: []primitive.ObjectID{}},
		{Name: "user", Permissions: []primitive.ObjectID{}},
	}

	
	permissionCollection := database.GetCollection("permissions")
	roleCollection := database.GetCollection("roles")

	
	for _, permission := range permissions {
		var existingPermission models.Permission
		err := permissionCollection.FindOne(context.Background(), bson.M{"name": permission.Name}).Decode(&existingPermission)
		if err == mongo.ErrNoDocuments {
	
			_, err := permissionCollection.InsertOne(context.Background(), permission)
			if err != nil {
				log.Println("Error creating permission:", permission.Name)
			}
		}
	}

	
	var createdPermissions []models.Permission
	cursor, err := permissionCollection.Find(context.Background(), bson.M{})
	if err != nil {
		log.Fatal("Error fetching permissions:", err)
	}
	if err := cursor.All(context.Background(), &createdPermissions); err != nil {
		log.Fatal("Error decoding permissions:", err)
	}

	
	var adminPermissions []primitive.ObjectID
	var userPermissions []primitive.ObjectID

	
	for _, permission := range createdPermissions {
		adminPermissions = append(adminPermissions, permission.ID) 
		if permission.Name == "view_task" {
			userPermissions = append(userPermissions, permission.ID)
		}
	}

	
	for _, role := range roles {
		var existingRole models.Role
		err := roleCollection.FindOne(context.Background(), bson.M{"name": role.Name}).Decode(&existingRole)
		if err == mongo.ErrNoDocuments {
	
			var permissionsToAssign []primitive.ObjectID
			if role.Name == "admin" {
				permissionsToAssign = adminPermissions
			} else {
				permissionsToAssign = userPermissions
			}

	
			role.Permissions = permissionsToAssign
			_, err := roleCollection.InsertOne(context.Background(), role)
			if err != nil {
				log.Println("Error creating role:", role.Name)
			}
		}
	}
}
