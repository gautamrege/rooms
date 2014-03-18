package rooms

import (
	"encoding/csv"
	"fmt"
	"os"
  "net/http"
  "html/template"
)

type Room struct {
	Type    string
	Roomies []string
	Emails  []string
	Tickets []string
}

var rooms = make(map[string]*Room)
var	tickets = make(map[string]string)

func setup(filename string) (err error) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}

	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()

	// Convert array of arrays to Hash
	rooms = make(map[string]*Room)
	tickets = make(map[string]string)

	for _, value := range lines {

    if value[2] == "" {
      continue
    }

		// update tickets
		tickets[value[0]] = value[2]

		// Test for existance
		ticket, ok := rooms[value[2]]
		if ok {
			// Search by room number
			ticket.Roomies = append(ticket.Roomies, value[6])
			ticket.Emails = append(ticket.Emails, value[7])
			ticket.Tickets = append(ticket.Tickets, value[0])
		} else {
			// insert new key
			rooms[value[2]] = &Room{Type: value[5],
				Roomies: []string{value[6]},
				Emails:  []string{value[7]},
				Tickets: []string{value[0]},
			}
		}
	}
	return
}

//func raffle(rooms map[string]*Room, tickets map[string]string) {
func raffle() string {
	var id string
	for id = range tickets {
		break
	}

	ticket := rooms[tickets[id]]
	for i, k := range ticket.Tickets {
		if id == k {
			return ticket.Roomies[i]
		}
	}
  return "Oops .. again"
}

func roomies(id string, rooms map[string]*Room, tickets map[string]string) {
	fmt.Println(rooms[tickets[id]].Roomies)
	fmt.Println(rooms[tickets[id]].Emails)
}

func init() {
	err := setup("./rooms.csv")
	if err != nil {
		fmt.Println("WTF ..", err)
	}

  http.HandleFunc("/", handler)
  http.HandleFunc("/rooms", rooms_handler)
  http.HandleFunc("/raffle", raffle_handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
  layout, err := template.ParseFiles("layout.html")
  if err != nil {
    panic(err)
  }

  data :=  &Room{}
  layout.Execute(w, data)
}

func rooms_handler(w http.ResponseWriter, r *http.Request) {
  layout, err := template.ParseFiles("layout.html")
  if err != nil {
    panic(err)
  }

  data, ok := rooms[tickets[r.FormValue("ticket")]]
  if ok == false {
    data = &Room{Type: "No Ticket or Conference Pass or Guest"}
  }
  layout.Execute(w, data)
}

func raffle_handler(w http.ResponseWriter, r *http.Request) {
  layout, err := template.ParseFiles("raffle.html")
  if err != nil {
    panic(err)
  }

  winner := raffle()
  layout.Execute(w, winner)
}

