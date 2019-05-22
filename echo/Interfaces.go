package main

import (
	"fmt"
)

type SalaryCalculator interface {
	CalculateSalary() int
}

type Permanent struct {
	empId    int
	basicpay int
	pf       int
}

type Contract struct {
	empId    int
	basicpay int
}

//salary of perment employee is sum of basic pay and pf
func (p Permanent) CalculateSalary() int {
	return p.basicpay + p.pf
}

//salary of contract employee is the basic pay alone
func (c Contract) CalculateSalary() int {
	return c.basicpay
}

/*
total expense is calculated by iterating though the SalaryCalculator slice and summing
the salaries of the individual employees
*/

// Cach 1
func totalExpense1(s SalaryCalculator, v SalaryCalculator, m SalaryCalculator) {
	expense := 0

	expense = expense + v.CalculateSalary()+ s.CalculateSalary()+ m.CalculateSalary()

	fmt.Printf("Total Expense Per Month $%d\n", expense)
}

// Cach 2
func totalExpense2(s []SalaryCalculator) {
	expense := 0
	for _, v := range s {
		expense = expense + v.CalculateSalary()
	}
	fmt.Printf("Total Expense Per Month $%d\n", expense)
}


func main() {
	pemp1 := Permanent{1, 5000, 20}
	pemp2 := Permanent{2, 6000, 30}
	cemp1 := Contract{3, 3000}

	// Implement Interface SalaryCalculator

	// Cach 1
	var s, v, m SalaryCalculator
	s, v, m = pemp1, pemp2, cemp1
	totalExpense1(s, v, m)
	// 1.1
	employees2 := []SalaryCalculator{s, v, m}
	totalExpense2(employees2)

	// Cach 2
	employees := []SalaryCalculator{pemp1, pemp2, cemp1}

	totalExpense2(employees)
}