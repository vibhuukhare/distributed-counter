# Service Discovery

    We are not using a centralized server, so we will be achieving this by making the nodes find and register as peers dynamically.
        	- The first node starts without peers.
            - When starting a new node we will specify atleast one peer.
            - So the node registers with a peer, and updating its address list hence allowing the network to grow and achieving the        purpose of service discovery.
            -Created a hearbeat function which checks other nodes health if one becomes inactive it will remove that node.

# Distributed in-memory counter

    We are using two apis /increment and /count to increment and update the information about the nodes and its peers.
            - /increment--> update the counter and inform all the peers
            - /count--> syncs the latest count.

# Go concurrency

    Achieved concurency by using mutexes, goroutines to make the service discovery handle multiple requests without encountering any race issues.


# Testing

    # Start the first node
        go run main.go --port=8080 &

    # Start additional nodes with peer discovery
        PEERS="localhost:8080" go run main.go --port=8081 &
        PEERS="localhost:8080,localhost:8081" go run main.go --port=8082 &

    # Check connected peers
        curl http://localhost:8080/api/peers

    # Increment the counter on Node 1:
        curl -X POST http://localhost:8080/api/increment

    # Check the counter values
        curl http://localhost:8080/api/count
        curl http://localhost:8081/api/count
        curl http://localhost:8082/api/count

    # Send 100 concurrent requests for race condition
        for i in {1..100}; do 
            curl -X POST http://localhost:8080/api/increment & 
        done
        wait
