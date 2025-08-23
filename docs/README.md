# Backend Code Structure Documentation

## Overview

This backend is built using **Go** and the **Nakama game server framework**. The architecture follows a **modular domain-driven design** with clear separation of concerns, event-driven communication, and layered responsibility patterns.

## Architecture Principles

### 1. **Modular Domain Organization**
- Each business domain (e.g., `account`, `profile`) is encapsulated in its own module
- Modules are self-contained with their own initialization, handlers, and business logic
- Cross-domain communication happens through events rather than direct dependencies

### 2. **Layered Architecture**
Each module follows a consistent 3-layer pattern:
- **Handler Layer**: HTTP/RPC request handling and validation
- **Service Layer**: Business logic and domain operations  
- **Data Layer**: Models and data access patterns

### 3. **Event-Driven Communication**
- Modules communicate via events to maintain loose coupling
- Central event processor routes events to appropriate handlers
- Asynchronous processing prevents blocking operations

## Directory Structure

```
├── main.go                    # Application entry point and module initialization
├── modules/                   # Domain modules
│   ├── account/              # User account management domain
│   ├── common/               # Shared components across domains
│   └── utils/                # Utility functions
├── shared/                   # Cross-cutting concerns
├── docs/                     # Documentation
├── scripts/                  # Deployment and maintenance scripts
└── tests/                    # Test suites
```

## Core Components

### 1. Application Entry Point (`main.go`)

```go
func InitModule(ctx context.Context, logger runtime.Logger, db *sql.DB, 
               nk runtime.NakamaModule, initializer runtime.Initializer) error
```

**Purpose**: Central initialization point for all modules
**Responsibilities**:
- Initialize each domain module in sequence
- Provide shared dependencies (logger, database, Nakama runtime)
- Measure and log startup performance

### 2. Domain Modules (`modules/`)

Each domain module follows this structure:

#### Account Module (`modules/account/`)
```
account/
├── routes.go      # Module initialization and RPC/hook registration
├── handler.go     # Request handling and validation
├── service.go     # Business logic and domain operations
├── hooks.go       # Nakama lifecycle event handlers
└── constants.go   # Domain-specific constants
```

**Module Pattern**:
- **`routes.go`**: Single `InitModule()` function per domain that registers all endpoints
- **`handler.go`**: RPC handlers that validate input and call service layer
- **`service.go`**: Pure business logic with no HTTP concerns
- **`hooks.go`**: Nakama lifecycle hooks (before/after authentication, etc.)

**Example Flow**:
```
HTTP Request → Handler (validation) → Service (business logic) → Response
```

### 3. Common Modules (`modules/common/`)

Shared components used across multiple domains:

```
common/
├── models/           # Shared data structures
├── services/         # Shared business services
├── eventProcessor/   # Central event routing
├── eventEmitter/     # Event publication utilities
└── notifier/         # Notification system
```

#### Event Processing System
```go
func ProcessEvent(ctx context.Context, logger runtime.Logger, evt *api.Event) {
    switch evt.GetName() {
    case "account_updated":
        account.HandleAccountUpdatedEvent(ctx, logger, evt)
    case "profile_updated":
        // profile.HandleProfileUpdatedEvent(ctx, logger, evt)
    default:
        logger.Error("unrecognized event: %+v", evt)
    }
}
```

**Purpose**: Central event router that enables loose coupling between modules
**Benefits**:
- Modules don't need direct dependencies on each other
- Easy to add new event types and handlers
- Asynchronous processing prevents blocking

### 4. Utilities (`modules/utils/`)

Safe execution patterns and helper functions:

```go
func SpawnSafe(parent context.Context, logger runtime.Logger, fn func(ctx context.Context)) {
    // Provides timeout, panic recovery, and proper context handling
}
```

**Purpose**: Reliable goroutine execution with error handling

### 5. Shared Components (`shared/`)

Cross-cutting concerns like error definitions:

```go
var (
    ErrBadInput           = runtime.NewError("input contained invalid data", INVALID_ARGUMENT)
    ErrInternalError      = runtime.NewError("internal server error", INTERNAL)
    ErrNoGuildFound       = runtime.NewError("guild not found", NOT_FOUND)
)
```

**Purpose**: Standardized error handling across all modules

## Data Flow Patterns

### 1. **Request-Response Flow**
```
Client Request
    ↓
Handler (validation + auth)
    ↓
Service (business logic)
    ↓
Nakama APIs / Database
    ↓
Response Models
    ↓
JSON Response
```

### 2. **Event-Driven Flow**
```
Service Operation
    ↓
Event Emission
    ↓
Event Processor (routing)
    ↓
Domain Event Handlers
    ↓
Side Effects (notifications, etc.)
```

### 3. **Hook-Based Flow**
```
Nakama Lifecycle Event
    ↓
Before/After Hooks
    ↓
Custom Business Logic
    ↓
Continue/Modify Default Behavior
```

## Key Design Decisions

### 1. **Why Modular Architecture?**
- **Scalability**: Teams can work on different domains independently
- **Maintainability**: Clear boundaries reduce complexity
- **Testability**: Each module can be tested in isolation
- **Reusability**: Common functionality is centralized

### 2. **Why Event-Driven Communication?**
- **Loose Coupling**: Modules don't directly depend on each other
- **Extensibility**: Easy to add new event handlers without modifying existing code
- **Asynchronous Processing**: Non-blocking operations improve performance

### 3. **Why Layered Within Modules?**
- **Separation of Concerns**: Each layer has a single responsibility
- **Testability**: Business logic is isolated from HTTP concerns
- **Maintainability**: Clear structure makes code easier to understand

### 4. **Why Nakama Framework?**
- **Real-time Features**: Built-in support for multiplayer, chat, notifications
- **Authentication**: OAuth, device auth, custom auth out of the box
- **Storage**: User accounts, data storage, leaderboards
- **Scalability**: Designed for game workloads

## Module Registration Pattern

Each module follows this initialization pattern:

```go
func InitModule(ctx context.Context, logger runtime.Logger, db *sql.DB, 
               nk runtime.NakamaModule, initializer runtime.Initializer) error {
    logger.Info("Initializing [Domain] domain...")
    
    // Register RPC endpoints
    utils.SpawnSafe(ctx, logger, func(ctx context.Context) {
        initializer.RegisterRpc("endpoint_name", HandlerFunction)
    })
    
    // Register lifecycle hooks
    utils.SpawnSafe(ctx, logger, func(ctx context.Context) {
        initializer.RegisterBeforeAuthenticateDevice(BeforeAuthHook)
        initializer.RegisterAfterAuthenticateDevice(AfterAuthHook)
    })
    
    logger.Info("[Domain] domain initialized")
    return nil
}
```

**Benefits**:
- Consistent initialization across modules
- Safe concurrent registration with timeout and panic recovery
- Clear logging for debugging initialization issues

## Error Handling Strategy

### 1. **Standardized Error Codes**
Using gRPC-style error codes for consistency:
```go
const (
    INVALID_ARGUMENT    = 3
    NOT_FOUND          = 5
    PERMISSION_DENIED  = 7
    INTERNAL           = 13
    UNAUTHENTICATED    = 16
)
```

### 2. **Structured Error Responses**
```go
return "", runtime.NewError("Authentication required", UNAUTHENTICATED)
```

### 3. **Panic Recovery**
All goroutines use `SpawnSafe()` for automatic panic recovery and logging.

## Extension Points

### Adding New Domains
1. Create new module directory under `modules/`
2. Implement the standard files: `routes.go`, `handler.go`, `service.go`
3. Add initialization call in `main.go`
4. Define any new events in the event processor

### Adding New Events
1. Add event case in `eventProcessor/event_processor.go`
2. Implement handler in the appropriate domain module
3. Emit events from service layer using `eventEmitter`

### Adding New Shared Services
1. Create service in `modules/common/services/`
2. Define models in `modules/common/models/`
3. Import and use across multiple modules

This architecture provides a solid foundation for a scalable, maintainable game backend that can grow with your requirements while maintaining clean separation of concerns and testability.
