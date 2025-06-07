package models

// Product represents a product in the system
// swagger:model
// @Description Product entity
// @name Product
// @Param id body string false "Product ID"
// @Param name body string true "Product Name"
// @Param description body string false "Description"
// @Param price body number true "Price"
// @Param image_url body string false "Image URL"
// @Param partner_id body string true "Partner ID"
// @Param stock body int true "Stock"
// @Param created_at body string false "Created At"
// @Param updated_at body string false "Updated At"
type Product struct {
	ID          string  `db:"id" json:"id"`
	Name        string  `db:"name" json:"name"`
	Description string  `db:"description" json:"description"`
	Price       float64 `db:"price" json:"price"`
	ImageURL    string  `db:"image_url" json:"image_url"`
	PartnerID   string  `db:"partner_id" json:"partner_id"`
	Stock       int     `db:"stock" json:"stock"`
	CreatedAt   string  `db:"created_at" json:"created_at"`
	UpdatedAt   string  `db:"updated_at" json:"updated_at"`
}

// ProductCreateRequest for creating a product
// swagger:model
// @Description Product create request
// @name ProductCreateRequest
type ProductCreateRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	ImageURL    string  `json:"image_url"`
	Stock       int     `json:"stock"`
}

// ProductUpdateRequest for updating a product
// swagger:model
// @Description Product update request
// @name ProductUpdateRequest
type ProductUpdateRequest struct {
	Name        *string  `json:"name,omitempty"`
	Description *string  `json:"description,omitempty"`
	Price       *float64 `json:"price,omitempty"`
	ImageURL    *string  `json:"image_url,omitempty"`
	Stock       *int     `json:"stock,omitempty"`
}
