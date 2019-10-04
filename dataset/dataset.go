package dataset

import (
	"encoding/json"
	"os"
)

var (
	PrimeCities = []int{122, 124, 462, 130, 123, 733, 141, 551, 134, 126}
	AddonCities = []int{
		// Bangalore, Hyderabad, Mumbai, Pune, Chennai, Delhi, Coimbatore,
		122, 124, 462, 130, 123, 733, 141,
		// Ahmedabad, Visakhapatnam, Tirupati, Jaipur, Goa, Pondicherry,
		551, 248, 71756, 807, 210, 403, 233,
		// Thanjavur, Kumbakonam, Jodhpur, Mysore, Chandigarh, Udupi, Kozhikode,
		66007, 427, 1169, 129, 833, 154, 74661,
		// Tiruchendur, Agra, Ooty, Amritsar, Chidambaram, Kodaikanal, Varanasi,
		663, 1290, 254, 759, 231, 227, 70429,
		// Haridwar, Digha, Shimla, Ajmer, Palani, Rameswaram, Ujjain
		802, 74706, 1285, 808, 501, 496, 1001,
	}
	PrimeBusOperators = []string{"bd8995a7-db41-4d0c-9e7f-d17ce9d896c0", "b4db17a8-027f-4598-8ff6-fbb4acb4b868", "05e8004b-3384-44d5-bc7b-a05b7176852d", "26a16e27-ac8a-4c1a-b8a4-4d1975ab337a", "3cae73e4-ac3d-4244-9e6a-db4050c84c27", "46fc7fdf-fc99-445b-b9ba-b3554122273c", "f456aae6-e622-4cd3-b079-7670cf0e50c4", "901d7916-d59d-4e4a-9297-e19fe0be42d2", "a490f88d-5315-4fe4-ae27-d92eb9cda599", "bed55e45-bd79-42f1-9b72-8e836eb7f4a0", "123b68f5-72f1-444e-95c5-a18035f4430c", "d8e14b87-3b53-420e-8272-76a068339f5f", "fed5cf89-7364-44df-9087-a7459e55d0e4", "c75b3315-0806-45b2-96fd-3a9ba17f3e3e", "113d927a-4013-4ee4-9523-9e86e92c706f", "b2755dd3-ec22-4278-8e6e-b1d505816bdb", "83b76f45-5f8d-4412-8361-5c896246bcd3", "29142750-4389-4e56-96d6-32ee8cac0b51", "75b4db11-7975-442c-864f-e8a27c730296", "93e9aa5d-0fbf-4da1-977a-72f680de1125", "8b62ff63-9311-4936-8e47-0af53d61ee62", "4a13de4c-2c10-4238-a102-5b26faab06a4", "67d4b760-6913-4ca5-86dc-44f4ee2ba9ef", "6dc242a8-5b7a-47cf-97d2-817bc082429e", "3bc8a394-eb66-4989-92e8-3e6a869e5d4f", "95311699-0aec-49c7-8148-c3b260e7ec61", "8e196bc1-a30a-4cb0-aa05-8726a7fb0844", "0fe0fe8b-f553-4c1c-b0de-d6d504903ce9", "ec43eb56-e953-4b0b-809d-6564b72fbca2", "0ccafa3e-b2da-445c-a10d-fb096c2d6cf2", "8c35fe2f-2668-445d-b0a9-44f635cf2187", "13be3ab7-898a-4db2-84b2-1f24f405730b", "19fb9174-1b93-41a0-89cb-329e95f7702c", "795a4427-5646-42ad-9185-bf1d10c902b5", "402896a3-36c5-42f5-8748-9e618a2455d4", "77139ec5-ef2f-4f41-aabb-1cf33c7eb98e", "ae0d4b59-9a30-4945-a96e-8e10e998c857", "7937c5b4-7d24-43c0-a864-0362c275437f", "0b5cc40d-e4ce-4ece-95f9-bb81ac2da34a", "a1c7a11b-c7da-470e-adf6-e29eea1a3d82", "96399ef0-9cb0-4c63-807d-ef83576b0999", "3dcbdfb8-6bf2-4bda-aa69-7a1472512990", "ca929e68-1687-437a-85e5-e52216e34f97", "b59910a5-6e9b-4292-a053-a36c6118420e", "3583fdc9-385c-4cc5-b3ee-59406b12fa8a", "91b5d90e-693d-4eeb-a7f3-7d9bc40f5848", "d587d69a-1118-45d3-8f1d-9641ea5a74a4", "c0b40df5-d736-4b36-8701-74d78b818fb7", "7ef413c5-b76e-4fb1-858d-ce9a906017ba", "d712207f-f536-4893-8546-8d2011f88c18"}
	Channels          = []string{"WEB_DIRECT", "MOBILE_APP", "MOBILE_WEB"}

	Cities       []int
	BusOperators []string
	Addons       []string
	CityBusPts   map[int][]int
)

func init() {
	f, err := os.Open("./dataset/cities.json")
	check(err)
	check(json.NewDecoder(f).Decode(&Cities))

	f, err = os.Open("./dataset/city-bus-pts.json")
	check(err)
	check(json.NewDecoder(f).Decode(&CityBusPts))

	f, err = os.Open("./dataset/bus_operators.json")
	check(err)
	check(json.NewDecoder(f).Decode(&BusOperators))

	f, err = os.Open("./dataset/addon_uuids.json")
	check(err)
	check(json.NewDecoder(f).Decode(&Addons))
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
