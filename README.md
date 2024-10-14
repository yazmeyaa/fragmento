# Fragmento

`fragmento` is a Go package designed for the fragmentation and assembly of data transmitted over the network. It utilizes a fragmentation mechanism to break large messages into smaller parts that can be sent over UDP or other data transmission methods.

## Installation

To install the package, use:

```bash
go get github.com/yazmeyaa/fragmento
```

## UDP Request Structure

| Field         | Size (bytes) | Description                             |
| ------------- | ------------ | --------------------------------------- |
| ID            | 4            | Unique identifier for the request       |
| Packed Fields | 1            | Packed fields                           |
| Index         | 2            | Index of the current fragment (0-based) |
| Total         | 2            | Total number of fragments               |
| Size          | 2            | Size of the payload (data)              |
| Payload       | Var.         | Payload (data)                          |
| Checksum      | 4            | Checksum for verifying data integrity   |

### Packed Fields Description
- **0:** Fragmented Flag (1 bit)
- **1-7:** RESERVED (7 bits)

### Fragmented Flag
The Fragmented Flag indicates whether the packet is a fragment of a larger message. If this flag is set to `1`, it signifies that the current packet is part of a fragmented transmission, allowing the receiver to identify and process it correctly within the context of the complete message.

## Core Entities

### Fragment

The `Fragment` struct represents a data fragment containing a header, size, payload, and checksum.

#### Methods

- **`Serialize() []byte`**  
  Serializes the fragment into a byte slice for transmission over the network.

### Header

The `Header` struct represents the header of a fragment, which includes an identifier, fragmentation flag, fragment index, and total number of fragments.

#### Methods

- **`Serialize() []byte`**  
  Serializes the header into a byte slice.

## Main Functions

### `FragmentData(id uint32, data []byte) []Fragment`

Function for fragmenting data. Takes an identifier and a byte array, returning a slice of fragments.

### `FromFragments(frags []Fragment) []byte`

Function that assembles data from a slice of fragments and returns the original byte array.

### `Deserialize(data []byte) (*Fragment, error)`

Function to deserialize a byte slice into a fragment object. Returns a pointer to the fragment and an error (if occurred).

## Example Usage

```go
package main

import (
    "fmt"
    "github.com/yazmeyaa/fragmento"
)

func main() {
    id := uint32(1)
    data := []byte("Hello, this is a message that needs to be fragmented.")

    // Fragmenting data
    fragments := fragmento.FragmentData(id, data)

    // Serializing fragments
    for _, frag := range fragments {
        serialized := frag.Serialize()
        fmt.Println("Serialized fragment:", serialized)
        
        // Deserializing fragment
        deserializedFrag, err := fragmento.Deserialize(serialized)
        if err != nil {
            fmt.Println("Deserialization error:", err)
            continue
        }
        
        // Checking data
        fmt.Println("Deserialized fragment:", deserializedFrag)
    }

    // Assembling data from fragments
    reconstructedData := fragmento.FromFragments(fragments)
    fmt.Println("Reconstructed data:", string(reconstructedData))
}
```

## License

This project is licensed under the [GNU Affero General Public License v3.0](LICENSE).

## Contributing

Contributions are welcome! Please feel free to fork the project, submit issues, or open pull requests.