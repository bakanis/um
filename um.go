package um

import "fmt"

var managers = make(map[string]UserManager)

// Register registers a user manager
func Register(name string, manager UserManager) {
	if name == "" {
		panic("um: Manager name cannot be empty")
	}
	if _, exists := managers[name]; exists == true {
		panic(fmt.Sprintf("um: Manager named '%s' already exists", name))
	}
	managers[name] = manager
}

// Open opens a user manager session
func Open(name, dns string) (UserManager, error) {
	manager, exists := managers[name]
	if exists == false {
		panic(fmt.Sprintf("um: Manager named '%s' does not exist", name))
	}
	err := manager.Setup(dns) // setup manager
	if err != nil {
		return nil, err
	}
	return manager, nil
}
