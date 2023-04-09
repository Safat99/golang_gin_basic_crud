package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// model for courses -file
type Course struct {
	CourseId    string  `json:"courseid"`
	CourseName  string  `json:"coursename"`
	CoursePrice int     `json:"price"`
	Author      *Author `json:"author"`
}

type Author struct {
	Fullname string `json:"fullname`
	Website  string `json:"website"`
}

// fake DB
var courses []Course

// middleware, helper -- file
func (c *Course) IsEmpty() bool {
	// return c.CourseId == "" && c.CourseName == ""
	return c.CourseName == ""
}

func main() {

	r := gin.Default()
	r.GET("/", serveHome)
	r.GET("/course/all", getAllCourses)
	r.GET("/course/:id", getSingleCourse)
	r.POST(("/course/add"), createOneCourse)
	r.PUT("/course/update/:id", updateOneCourse)
	r.DELETE("/course/delete/:id", deleteOneCourse)
	r.DELETE("/course/all", deleteAllCourse)

	r.Run(":8080")
}

//controllers -- different file

//serve home route

func serveHome(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Welcome to Home Page",
	})
	// c.Header("Content-Type", "text/html")
	// html := "<h1>Hello World!</h1>"
	// fmt.Fprint(c.Writer, html)
}

func getAllCourses(c *gin.Context) {
	fmt.Println("Get all courses")
	c.Header("Content-Type", "application/json")
	json.NewEncoder(c.Writer).Encode(courses)
}

func getSingleCourse(c *gin.Context) {
	fmt.Println("------------Getting One course-------")
	c.Header("Content-Type", "application/json")

	// grab id from request
	id := c.Param("id")

	// loop through courses and find matching id and return the response
	for _, course := range courses {
		if course.CourseId == id {
			json.NewEncoder(c.Writer).Encode(course) // line for mux
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{
		"error": "User not found",
	})

	// json.NewEncoder(c.Writer).Encode("No Course found with the given id")
	// return
}

func createOneCourse(c *gin.Context) {

	fmt.Println("------------Create One course-------")
	c.Header("Content-Type", "application/json")

	// what about {} -->> user is sending {}

	var course Course
	fmt.Println(course)
	decoder := json.NewDecoder(c.Request.Body) // makes a decoder struct based on our rq body

	err := decoder.Decode(&course)
	//Pointer pass korteso function e,
	//oi funtion ta json ke object e convert korbe ei address er type er upor base kore
	fmt.Println(course)
	if err != nil {
		fmt.Println(course)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	course.CourseId = strconv.Itoa(rand.Intn(100) + 1)

	courses = append(courses, course)

	c.JSON(http.StatusCreated, course)

}

func updateOneCourse(c *gin.Context) {
	fmt.Println("Updating one course")
	c.Header("Content-Type", "application/json")

	// grab the id
	id := c.Param("id")

	var isFound bool = false

	// loop, find matching id, remove, add with given id
	for index, course := range courses {
		if course.CourseId == id {
			isFound = true
			courses = append(courses[:index], courses[index+1:]...)

			var course Course
			decoder := json.NewDecoder(c.Request.Body)
			if err := decoder.Decode(&course); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
				return
			}

			course.CourseId = id
			courses = append(courses, course)

			c.JSON(http.StatusOK, course)
			return
		}
	}

	if !isFound {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "no course found against " + id + " id",
		})
	}

}

func deleteOneCourse(c *gin.Context) {
	fmt.Println("Deleting one course")
	c.Header("Content-Type", "application/json")

	// grab the id
	id := c.Param("id")
	var isFound bool = false

	for index, course := range courses {
		if course.CourseId == id {
			isFound = true
			courses = append(courses[:index], courses[index+1:]...)

			c.JSON(http.StatusOK, gin.H{
				"message": "course id " + id + " removed successfully!!!",
			})
			return
		}
	}

	if !isFound {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "no course found against " + id + " id",
		})
		return
	}

}

func deleteAllCourse(c *gin.Context) {
	fmt.Println("----------deleting all courses-------")
	c.Header("Content-Type", "application/json")

	courses = nil

	c.JSON(http.StatusOK, gin.H{
		"message": "all courses are removed!!!",
	})

}
