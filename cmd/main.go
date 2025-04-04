package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"screen-therapy-backend/config"
	"screen-therapy-backend/routes"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize Firebase & Firestore
	config.InitFirebase()
	defer config.Client.Close()

	// Create a new router
	r := mux.NewRouter()
	routes.RegisterRoutes(r)

	// Print IP address
	printLocalIP()

	// Start server
	fmt.Println("🚀 Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func printLocalIP() {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println("Could not get local IP address:", err)
		return
	}

	for _, addr := range addrs {
		// Only show IPv4 and skip loopback
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			ip := ipNet.IP.To4()
			if ip != nil {
				fmt.Printf("🌐 Access from local network: http://%s:8080\n", ip.String())
				return
			}
		}
	}
}
