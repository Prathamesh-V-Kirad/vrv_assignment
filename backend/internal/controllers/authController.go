package controllers

import (
	"context"
	"backend/internal/models"
	"backend/internal/database"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"go.mongodb.org/mongo-driver/bson"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


func Register(c *fiber.Ctx) error {
	var data map[string]string
	
	if err := c.BodyParser(&data); err != nil {
		return err
	}
	
	password,_ := bcrypt.GenerateFromPassword([]byte(data["password"]),14)
	user := models.User{
		Name : data["name"],
		Email : data["email"],
		Password: password,
	}
	collection := database.GetCollection("users")

	// Insert the user into the database
	insertResult, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	user.ID = insertResult.InsertedID.(primitive.ObjectID)
	return c.Status(fiber.StatusCreated).JSON(user)
}

func Login(c *fiber.Ctx) error {
    var data map[string]string

    if err := c.BodyParser(&data); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid request body",
        })
    }

    email := data["email"]
    password := data["password"]
    if email == "" || password == "" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Email and password are required",
        })
    }

    collection := database.GetCollection("users")
    var user models.User
    err := collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": "Invalid email or password",
        })
    }

    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": "Invalid email or password",
        })
    }

    expirationTime := time.Now().Add(1 * time.Hour)
    claims := &models.CustomClaims{
        Role:   user.RoleID.Hex(),
        RegisteredClaims: jwt.RegisteredClaims{
            Issuer:    user.ID.Hex(),
            ExpiresAt: jwt.NewNumericDate(expirationTime),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            ID:        generateJTI(),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    secretKey := os.Getenv("JWT_SECRET_KEY")
    if secretKey == "" {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "JWT Secret key is not set",
        })
    }

    signedToken, err := token.SignedString([]byte(secretKey))
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to sign the token",
        })
    }

    cookie := fiber.Cookie{
        Name:     "jwt",
        Value:    signedToken,
        Expires:  expirationTime,
        HTTPOnly: true,
        Secure:   false,
    }

    c.Cookie(&cookie)
    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "message": "Login successful",
        "token":   signedToken,
    })
}

func generateJTI() string {
    return primitive.NewObjectID().Hex()
}

func User(c *fiber.Ctx) error {
    cookie := c.Cookies("jwt")
    if cookie == "" {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": "Missing token",
        })
    }

    secretKey := os.Getenv("JWT_SECRET_KEY")
    if secretKey == "" {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "JWT secret key is not set",
        })
    }

    
    token, err := jwt.ParseWithClaims(cookie, &models.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid signing method")
        }
        return []byte(secretKey), nil
    })

    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": "Invalid or expired token",
        })
    }

    
    if claims, ok := token.Claims.(*models.CustomClaims); ok && token.Valid {    
        userID, err := primitive.ObjectIDFromHex(claims.Issuer)
        if err != nil {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error": "Invalid user ID in token",
            })
        }

        collection := database.GetCollection("users")
        var user models.User
        err = collection.FindOne(context.Background(), bson.M{"_id": userID}).Decode(&user)
        if err != nil {
            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
                "error": "User not found",
            })
        }

        
        return c.JSON(user)
    }

    return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
        "error": "Invalid token",
    })
}

func Logout(c *fiber.Ctx) error {
    cookie := fiber.Cookie{
        Name: "jwt",
        Value: "",
        Expires: time.Now().Add(-time.Hour),
        HTTPOnly: true,
    }
    c.Cookie(&cookie)
    return c.JSON(fiber.Map{
        "message": "success",
    })


}