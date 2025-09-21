package requestresponse

type ServicesReq struct {
	PetSitting  bool `json:"pet_sitting"`
	DogWalking  bool `json:"dog_walking"`
	PetDayCare  bool `json:"pet_day_care"`
	PetGrooming bool `json:"pet_grooming"`
	PetTraining bool `json:"pet_training"`
	PetMassage  bool `json:"pet_massage"`
}

type UpdateServicesReq struct {
	PetSitting  *bool `json:"pet_sitting" validate:"required"`
	DogWalking  *bool `json:"dog_walking"`
	PetDayCare  *bool `json:"pet_day_care"`
	PetGrooming *bool `json:"pet_grooming"`
	PetTraining *bool `json:"pet_training"`
	PetMassage  *bool `json:"pet_massage"`
}
