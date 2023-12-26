package handlers

import (
	"errors"
	"mentalartsapi_hw/models"

	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	db *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
	h := Handler{db}
	return &h
}

func (h *Handler) InitRoutes(r *gin.Engine) {
	// User routes
	userRoutes := r.Group("/users")
	{
		userRoutes.POST("/", h.createUser)
		userRoutes.GET("/", h.getUsers)
		userRoutes.GET("/:id", h.getUser)
		userRoutes.GET("/:id/posts", h.getUserPosts)
		userRoutes.DELETE("/:id", h.deleteUser)
	}

	// Post routes
	postRoutes := r.Group("/posts")
	{
		postRoutes.POST("/", h.createPost)
		postRoutes.GET("/:id", h.getPost)
		postRoutes.PUT("/:id", h.updatePost)
		userRoutes.GET("/:id/comments", h.getPostComments)
		postRoutes.DELETE("/:id", h.deletePost)
	}

	// Comment routes
	commentRoutes := r.Group("/comments")
	{
		commentRoutes.POST("/", h.createComment)
		commentRoutes.GET("/:id", h.getComment)
		commentRoutes.PUT("/:id", h.updateComment)
		commentRoutes.DELETE("/:id", h.deleteComment)
	}
}

// @Summary Create a new user
// @Description Create a new user
// @ID create-user
// @Accept json
// @Produce json
// @Param user body User true "User object to be created"
// @Success 200 {object} User
// @Router /users [post]
func (h *Handler) createUser(c *gin.Context) {
	var newUser models.User
	if err := c.BindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	newUser.ID = 0

	if err := h.db.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, newUser)
}

// @Summary Get all users
// @Description Get a list of all users
// @ID get-users
// @Produce json
// @Success 200 {array} User
// @Router /users [get]
func (h *Handler) getUsers(c *gin.Context) {
	var users []models.User

	if err := h.db.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, users)
}

// @Summary Get user by ID
// @Description Get a user by ID
// @ID get-user
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} User
// @Router /users/{id} [get]
func (h *Handler) getUser(c *gin.Context) {
	idParam := c.Param("id")
	var author models.User

	id, err := strconv.Atoi(idParam)

	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if err := h.db.First(&author, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, err.Error())
			return
		}
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, author)
}

// @Summary Get posts by user ID
// @Description Get posts created by a specific user
// @ID get-user-posts
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {array} Post
// @Router /users/{id}/posts [get]
func (h *Handler) getUserPosts(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, "Invalid user ID")
		return
	}

	var user models.User
	if err := h.db.Preload("Posts").First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, "User not found")
			return
		}
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, user.Posts)
}

// @Summary Delete user by ID
// @Description Delete a user and associated data by ID
// @ID delete-user
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {string} string
// @Router /users/{id} [delete]
func (h *Handler) deleteUser(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, "Invalid user ID")
		return
	}

	// Fetch the user
	var user models.User
	if err := h.db.Preload("Posts.Comments").First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, "User not found")
			return
		}
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	// Delete associated comments
	for _, post := range user.Posts {
		for _, comment := range post.Comments {
			if err := h.db.Delete(&comment).Error; err != nil {
				c.JSON(http.StatusInternalServerError, err.Error())
				return
			}
		}
	}

	// Delete associated posts
	for _, post := range user.Posts {
		if err := h.db.Delete(&post).Error; err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
	}

	// Delete the user
	if err := h.db.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, "User and associated data deleted successfully")
}

// @Summary Create a new post
// @Description Create a new post
// @ID create-post
// @Accept json
// @Produce json
// @Param userID path int true "User ID"
// @Param post body Post true "Post object to be created"
// @Success 200 {object} Post
// @Router /posts [post]
func (h *Handler) createPost(c *gin.Context) {
	var newPost models.Post
	if err := c.BindJSON(&newPost); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	userID := newPost.UserID
	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusBadRequest, "Invalid UserID")
		return
	}

	newPost.ID = 0

	if err := h.db.Create(&newPost).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, newPost)
}

// @Summary Get post by ID
// @Description Get a post by ID
// @ID get-post
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {object} Post
// @Router /posts/{id} [get]
func (h *Handler) getPost(c *gin.Context) {
	postIDParam := c.Param("id")
	var post models.Post

	postID, err := strconv.Atoi(postIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Invalid post ID")
		return
	}

	if err := h.db.Preload("Comments").First(&post, postID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, "Post not found")
			return
		}

		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, post)
}

// @Summary Get comments for a post
// @Description Get comments for a specific post
// @ID get-post-comments
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {array} Comment
// @Router /posts/{id}/comments [get]
func (h *Handler) getPostComments(c *gin.Context) {
	postIDParam := c.Param("id")
	postID, err := strconv.Atoi(postIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Invalid post ID")
		return
	}

	var postWithComments models.Post

	if err := h.db.Preload("Comments").First(&postWithComments, postID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, "Post not found")
			return
		}
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, postWithComments)
}

// @Summary Update a post by ID
// @Description Update a post by ID
// @ID update-post
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Param post body Post true "Updated post object"
// @Success 200 {object} Post
// @Router /posts/{id} [put]
func (h *Handler) updatePost(c *gin.Context) {
	postIDParam := c.Param("id")
	postID, err := strconv.Atoi(postIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Invalid post ID")
		return
	}

	var updatedPost models.Post

	if err := h.db.First(&updatedPost, postID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, "Post not found")
			return
		}
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	if err := c.BindJSON(&updatedPost); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if err := h.db.Save(&updatedPost).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, updatedPost)
}

// @Summary Delete post by ID
// @Description Delete a post by ID
// @ID delete-post
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {string} string
// @Router /posts/{id} [delete]
func (h *Handler) deletePost(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, "Invalid post ID")
		return
	}

	var post models.Post
	if err := h.db.First(&post, postID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, "Post not found")
			return
		}
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	if err := h.db.Delete(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, "Post deleted successfully")
}

// @Summary Create a new comment
// @Description Create a new comment
// @ID create-comment
// @Accept json
// @Produce json
// @Param postID path int true "Post ID"
// @Param comment body Comment true "Comment object to be created"
// @Success 200 {object} Comment
// @Router /comments [post]
func (h *Handler) createComment(c *gin.Context) {
	var newComment models.Comment
	if err := c.BindJSON(&newComment); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	newComment.ID = 0

	if err := h.db.Create(&newComment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, newComment)
}

// @Summary Get comment by ID
// @Description Get a comment by ID
// @ID get-comment
// @Produce json
// @Param id path int true "Comment ID"
// @Success 200 {object} Comment
// @Router /comments/{id} [get]
func (h *Handler) getComment(c *gin.Context) {
	commentIDParam := c.Param("id")
	commentID, err := strconv.Atoi(commentIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Invalid comment ID")
		return
	}

	var comment models.Comment

	if err := h.db.First(&comment, commentID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, "Comment not found")
			return
		}
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, comment)
}

// @Summary Update comment by ID
// @Description Update a comment by ID
// @ID update-comment
// @Accept json
// @Produce json
// @Param id path int true "Comment ID"
// @Param input body Comment true "Comment object to be updated"
// @Success 200 {object} Comment
// @Router /comments/{id} [put]
func (h *Handler) updateComment(c *gin.Context) {
	commentIDParam := c.Param("id")
	commentID, err := strconv.Atoi(commentIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Invalid comment ID")
		return
	}

	var updatedComment models.Comment

	if err := c.BindJSON(&updatedComment); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	updatedComment.ID = uint(commentID)

	if err := h.db.Save(&updatedComment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedComment)
}

// @Summary Delete comment by ID
// @Description Delete a comment by ID
// @ID delete-comment
// @Produce json
// @Param id path int true "Comment ID"
// @Success 200 {object} Comment
// @Router /comments/{id} [delete]
func (h *Handler) deleteComment(c *gin.Context) {
	commentIDParam := c.Param("id")
	commentID, err := strconv.Atoi(commentIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Invalid comment ID")
		return
	}

	var deletedComment models.Comment

	if err := h.db.First(&deletedComment, commentID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, "Comment not found")
			return
		}
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	if err := h.db.Delete(&deletedComment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, deletedComment)
}
