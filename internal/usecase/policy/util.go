package policy

import (
	"fmt"

	policyRepo "github.com/aidapedia/jabberwock/internal/repository/policy"
)

// Will Generated Policy Like
// 1, http:GET, /api/v1/users
// 1, rpc:OrderService, GetOrder
func StdPolicy(policy policyRepo.Policy) []interface{} {
	return []interface{}{
		policy.Role,
		fmt.Sprintf("%s:%s", policy.Type, policy.Method),
		policy.Path,
	}
}
