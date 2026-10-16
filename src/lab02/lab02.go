package main

import "fmt"

type Pet interface {
	express() 
	getName() string
	getAge() int 
	getAgeInHumanYears() int
}

type Animal struct {
	name string 
	age int 

}


func (a *Animal) getName() string {
	return a.name
}

func (a *Animal) getAge() int {
	return a.age
}

func (a *Animal) happyBirthday() {
	a.age++
}

type Parrot struct {
	Animal
}

type Turtle struct {
	Animal
}

func (t *Turtle) express() {
	fmt.Println("   _____     ____")
	fmt.Println(" /      \\  |  o |")
	fmt.Println("|        |/ ___\\|")
	fmt.Println("|_________/")
	fmt.Println("|_|_| |_|_|")
}

func (t *Turtle) getAgeInHumanYears() int {
	return t.getAge() * 5 + 1
}

func (p *Parrot) express() {
	fmt.Println(p.getName(), ": \"Hello!\"")
	fmt.Println("    \\\\")
	fmt.Println("     (o>")
	fmt.Println(" \\_//)")
	fmt.Println("  \\_/_)")
	fmt.Println("   _|_")
}

func (p *Parrot) getAgeInHumanYears() int {
	return p.getAge() * 5
}

type Person struct {
	name string 
	pet Pet
}

func (p *Person) setPet(pet Pet) {
	p.pet = pet
}

func (p *Person) getPetInfo() {
	fmt.Println("--------------------")

	if p.pet == nil {
		fmt.Println(p.name, "has no pet yet")
		fmt.Println("--------------------")
		return
	}

	fmt.Println("Pet for", p.name, "is named", p.pet.getName())

    fmt.Print("Age: ")
	fmt.Print(p.pet.getAge())
	fmt.Print(" year(s), about ")
	fmt.Print(p.pet.getAgeInHumanYears())
	fmt.Println(" human years")

	fmt.Println("And it expresses itself like this:")
	p.pet.express()
	fmt.Println("--------------------")
}


func main() {
	jack := Person{name: "Jack", pet: nil}
	john := Person{name:"John", pet: nil}

	parrot := Parrot{Animal{name: "Kesha", age: 3}}
	turtle := Turtle{Animal{name: "Zavala", age: 4}}

	jack.setPet(&parrot)

	jack.getPetInfo()
	john.getPetInfo()

	john.setPet(&turtle)
	
	john.getPetInfo()
}