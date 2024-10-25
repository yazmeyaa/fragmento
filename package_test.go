package fragmento_test

import (
	"bytes"
	"encoding/binary"
	"hash/crc32"
	"math/rand"
	"testing"

	fragmento "github.com/yazmeyaa/fragmento"
)

func TestFragment(t *testing.T) {
	payload := []byte{10, 20, 30, 40, 50, 60}
	header := fragmento.NewHeader(1, false, 0, 0)
	fragment := fragmento.NewFragment(payload, header)
	result := fragment.Serialize()

	expectedChecksum := crc32.ChecksumIEEE(payload)

	checksum := make([]byte, 4)
	binary.LittleEndian.PutUint32(checksum, expectedChecksum)

	expectedBuffer := []byte{
		1, 0, 0, 0,
		0,
		0, 0,
		0, 0,
		6, 0,
		10, 20, 30, 40, 50, 60,
	}
	expectedBuffer = append(expectedBuffer, checksum...)

	if !bytes.Equal(result, expectedBuffer) {
		t.Errorf("Buffers are not equals.\nExpected: %v\nGot: %v", expectedBuffer, result)
	}
}

func TestFramgentLargeData(t *testing.T) {
	id := rand.Uint32()
	payload := make([]byte, fragmento.MAX_PAYLOAD_SIZE+1000)
	for v := 0; v < len(payload); v++ {
		payload[v] = byte(v % 255)
	}
	frags := fragmento.FragmentData(id, payload)

	part1 := payload[:fragmento.MAX_PAYLOAD_SIZE]
	part2 := payload[fragmento.MAX_PAYLOAD_SIZE:]
	parts := []([]byte){part1, part2}

	for idx, x := range frags {
		eq := bytes.Equal(x.Payload, parts[idx])
		if !eq {
			t.Error("Not equal parts of data.")
			return
		}
	}

}

func TestDeserializtion(t *testing.T) {
	id := rand.Uint32()
	payload := make([]byte, 150)

	for v := 0; v < len(payload); v++ {
		payload[v] = byte(v % 255)
	}
	frags := fragmento.FragmentData(id, payload)

	gotBytes := frags[0].Serialize()
	frag, err := fragmento.Deserialize(gotBytes)
	if err != nil {
		t.Error(err)
	}

	expectedChecksum := crc32.ChecksumIEEE(payload)
	expectedId := id
	expectedLenght := len(payload)

	if expectedChecksum != frag.Checksum {
		t.Error("Wrong checksum after deserialization")
		return
	}
	if expectedId != frag.Header.ID() {
		t.Error("Wrong ID after deserialization")
		return
	}
	if expectedLenght != int(frag.Size) {
		t.Error("Wrong checksum payload size after deserialization")
		return
	}
	if !bytes.Equal(frag.Payload, payload) {
		t.Error("Wrong payload after deserialization")
		return
	}
}

func TestFromFragments(t *testing.T) {
	text := `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Pellentesque odio mi, interdum eget malesuada ut, scelerisque porttitor lorem. Vestibulum nec dui tristique, aliquet elit non, tempus odio. Vivamus ultricies tempor tristique. Nullam quis enim porta leo volutpat fringilla. Etiam bibendum, sapien at tincidunt imperdiet, mi eros vulputate est, ac tristique diam ex eget urna. Aenean blandit consequat lorem ac ultrices. Maecenas molestie mauris justo, at dapibus erat sollicitudin ultricies. Etiam sodales leo ut lorem iaculis, sed molestie urna iaculis. Sed fermentum leo in tincidunt interdum.
Sed malesuada, turpis at ultrices varius, magna libero molestie nisl, sit amet porttitor nibh lectus at metus. Donec sit amet pellentesque velit. Cras quis feugiat ligula. Nam at sem nisi. Maecenas a urna molestie, pulvinar purus non, vulputate ipsum. Sed a velit euismod justo volutpat mattis. Mauris lorem libero, semper sit amet posuere sit amet, efficitur ut enim. Ut erat metus, hendrerit sed mattis pellentesque, venenatis sed dui. Maecenas iaculis semper est, ac aliquet eros lacinia a.
Ut egestas urna ut massa consectetur iaculis. Ut in elit bibendum, euismod libero egestas, malesuada tellus. Phasellus ut eros sed tellus feugiat mattis dapibus nec turpis. Aenean consequat porta tempus. Phasellus in nisi in eros volutpat eleifend sed sed diam. Vestibulum malesuada mattis ligula eu viverra. Quisque sit amet dictum orci. Suspendisse eget cursus nulla. Maecenas volutpat ornare consectetur. Quisque quam lorem, fermentum sit amet justo eu, faucibus condimentum dolor. Nunc ut neque et magna finibus vehicula sit amet non lectus. Nam in mollis urna, sodales bibendum urna. Duis vehicula erat nibh, id mattis massa porttitor a. Mauris tristique, est in feugiat tincidunt, ex diam lobortis diam, eget dapibus eros justo id ex. Praesent fermentum ex et maximus pulvinar. Curabitur nec quam in enim faucibus facilisis eu at dolor.
Proin sit amet elit eget metus consequat consectetur. Curabitur eget neque eget sapien sodales consequat. Sed consectetur leo ut mi rhoncus, a fermentum tortor euismod. Proin at leo ac velit tempor pharetra sed sed augue. Proin cursus sit amet elit congue fermentum. Donec porttitor volutpat tortor, vel ullamcorper metus ultricies in. In ut pretium turpis. Aliquam erat volutpat. Morbi sagittis nisl ut sagittis vulputate. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Proin varius venenatis sapien, in malesuada nibh auctor eu. Cras commodo nulla vitae posuere cursus. Curabitur sed lacinia tellus. Fusce laoreet at tellus eu pretium.
Morbi mattis tincidunt leo in vulputate. Sed vestibulum luctus ex ac ullamcorper. Fusce eu risus ut mi ornare semper a et sem. Morbi scelerisque, leo ut tincidunt tempus, tortor lorem posuere est, sit amet venenatis dolor odio id nibh. Aenean vitae mollis tellus, a tempus nunc. Etiam mattis odio mi, ut blandit tellus molestie quis. Pellentesque nec vehicula tortor, vitae volutpat nibh. Mauris congue nisl vitae dolor rutrum maximus. Nunc faucibus arcu turpis, maximus efficitur tortor suscipit et. Etiam congue lacus ut nunc vehicula rhoncus. Etiam a nibh eu magna porta luctus ac sit amet nunc.
Class aptent taciti sociosqu ad litora torquent per conubia nostra, per inceptos himenaeos. Donec a metus non turpis porttitor lobortis. Nulla ut ante magna. Etiam nec nisl congue, consequat tellus eget, scelerisque mi. Interdum et malesuada fames ac ante ipsum primis in faucibus. Donec vulputate, nisl sit amet faucibus bibendum, velit eros pharetra ligula, id blandit sapien turpis nec nulla. Praesent id placerat elit. Suspendisse faucibus laoreet orci, ut maximus dolor suscipit et. Morbi vel erat rhoncus, dignissim quam id, bibendum elit. Duis vestibulum quis dui vitae semper. Nulla consectetur non odio at facilisis.
Vestibulum ornare mollis feugiat. Etiam condimentum, ipsum et mollis mollis, orci tortor pharetra dolor, id fermentum leo felis a purus. Ut sed risus sed mauris consequat scelerisque a et eros. Etiam tristique neque ac libero sodales convallis. Nulla aliquam nibh neque, in semper est orci.`

	id := rand.Uint32()
	payload := []byte(text)
	fragments := fragmento.FragmentData(id, payload)

	data := fragmento.FromFragments(fragments)
	if !bytes.Equal(payload, data) {
		t.Error("Not equal data")
	}
}
