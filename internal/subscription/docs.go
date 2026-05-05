// Package subscription implements subscription management domain logic.
//
// It provides HTTP handlers, service layer, repository layer, and DTO/domain
// mapping for creating, updating, deleting, listing, and aggregating
// subscription data.
//
// The package follows a layered architecture:
//   - transport (HTTP handlers)
//   - service (business logic)
//   - repository (PostgreSQL persistence)
//   - mapping (DTO ↔ domain conversion)
package subscription
