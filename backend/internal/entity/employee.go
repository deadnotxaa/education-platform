// Package entity defines main entities for business logic (services), database mapping, and
// HTTP response objects if suitable. Each logic group entity in its own file.
package entity

type (
    // Role -.
    Role struct {
        ID   int    `json:"id"          example:"2"`
        Name string `json:"name"        example:"Administrator"`
    }

    // Employee -.
    Employee struct {
        ID     int `json:"id"           example:"1"`
        UserID int `json:"user_id"      example:"1"`
        RoleID int `json:"role_id"      example:"2"`
    }

    // Teacher -.
    Teacher struct {
        EmployeeID               int    `json:"employee_id"                 example:"1"`
        WorkPlace                string `json:"work_place"                  example:"Some BigTech Company"`
        OverallExperience        int    `json:"overall_experience"          example:"5"`
        SpecializationExperience int    `json:"specialization_experience"   example:"3"`
    }
)
