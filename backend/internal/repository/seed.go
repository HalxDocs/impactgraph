package repository

import (
	"fmt"

	"impactgraph/internal/models"
)

// seed populates the InMemoryStore with realistic data.
// Organizations: 3, Teams: 8, Developers: 20, Repositories: 30, Packages: 40, Services: 15
func seed(s *InMemoryStore) {
	// Organizations
	orgs := []models.Node{
		{ID: "org-acme", Type: "Organization", Label: "Organization", Name: "Acme Technologies", Properties: map[string]interface{}{"name": "Acme Technologies"}},
		{ID: "org-nova", Type: "Organization", Label: "Organization", Name: "Nova Systems", Properties: map[string]interface{}{"name": "Nova Systems"}},
		{ID: "org-orbit", Type: "Organization", Label: "Organization", Name: "Orbit Labs", Properties: map[string]interface{}{"name": "Orbit Labs"}},
	}
	for _, o := range orgs {
		s.addNode(o)
	}

	teams := []models.Node{
		{ID: "team-infra", Type: "Team", Label: "Team", Name: "Infrastructure", Properties: map[string]interface{}{"name": "Infrastructure", "org": "org-acme"}},
		{ID: "team-platform", Type: "Team", Label: "Team", Name: "Platform", Properties: map[string]interface{}{"name": "Platform", "org": "org-acme"}},
		{ID: "team-checkout", Type: "Team", Label: "Team", Name: "Checkout", Properties: map[string]interface{}{"name": "Checkout", "org": "org-acme"}},
		{ID: "team-payments", Type: "Team", Label: "Team", Name: "Payments", Properties: map[string]interface{}{"name": "Payments", "org": "org-nova"}},
		{ID: "team-notifications", Type: "Team", Label: "Team", Name: "Notifications", Properties: map[string]interface{}{"name": "Notifications", "org": "org-nova"}},
		{ID: "team-frontend", Type: "Team", Label: "Team", Name: "Frontend", Properties: map[string]interface{}{"name": "Frontend", "org": "org-orbit"}},
		{ID: "team-data", Type: "Team", Label: "Team", Name: "Data & Analytics", Properties: map[string]interface{}{"name": "Data & Analytics", "org": "org-orbit"}},
		{ID: "team-security", Type: "Team", Label: "Team", Name: "Security", Properties: map[string]interface{}{"name": "Security", "org": "org-acme"}},
	}
	for _, t := range teams {
		s.addNode(t)
	}

	developers := []struct {
		id, name, email, github string
	}{
		{"dev-1", "Alice Chen", "alice@acme.dev", "alicechen"},
		{"dev-2", "Bob Martinez", "bob@acme.dev", "bobmartinez"},
		{"dev-3", "Carol Kim", "carol@nova.dev", "carolkim"},
		{"dev-4", "David Okafor", "david@nova.dev", "davidoka"},
		{"dev-5", "Eva Patel", "eva@orbit.dev", "evapatel"},
		{"dev-6", "Frank Liu", "frank@acme.dev", "frankliu"},
		{"dev-7", "Grace Smith", "grace@acme.dev", "gracesmith"},
		{"dev-8", "Hassan Ali", "hassan@nova.dev", "hassanali"},
		{"dev-9", "Ivy Johnson", "ivy@orbit.dev", "ivyj"},
		{"dev-10", "Jack Wilson", "jack@acme.dev", "jackw"},
		{"dev-11", "Karen Brown", "karen@nova.dev", "karenb"},
		{"dev-12", "Leo Garcia", "leo@acme.dev", "leogarcia"},
		{"dev-13", "Maya Singh", "maya@orbit.dev", "mayasingh"},
		{"dev-14", "Noah Taylor", "noah@acme.dev", "noaht"},
		{"dev-15", "Olivia Davis", "olivia@nova.dev", "oliviad"},
		{"dev-16", "Paul Anderson", "paul@acme.dev", "paula"},
		{"dev-17", "Quinn Thomas", "quinn@orbit.dev", "quinnt"},
		{"dev-18", "Rachel Green", "rachel@acme.dev", "rachelg"},
		{"dev-19", "Sam White", "sam@nova.dev", "samwhite"},
		{"dev-20", "Tina Black", "tina@orbit.dev", "tinab"},
	}
	for _, d := range developers {
		s.addNode(models.Node{ID: d.id, Type: "Developer", Label: "Developer", Name: d.name, Properties: map[string]interface{}{"name": d.name, "email": d.email, "github": d.github}})
	}

	repos := []models.Node{
		{ID: "repo-checkout-api", Type: "Repository", Label: "Repository", Name: "checkout-api", Properties: map[string]interface{}{"name": "checkout-api", "url": "https://github.com/acme/checkout-api", "language": "Go", "visibility": "private"}},
		{ID: "repo-payment-service", Type: "Repository", Label: "Repository", Name: "payment-service", Properties: map[string]interface{}{"name": "payment-service", "url": "https://github.com/nova/payment-service", "language": "Go", "visibility": "private"}},
		{ID: "repo-notification-service", Type: "Repository", Label: "Repository", Name: "notification-service", Properties: map[string]interface{}{"name": "notification-service", "url": "https://github.com/nova/notification-service", "language": "TypeScript", "visibility": "private"}},
		{ID: "repo-user-service", Type: "Repository", Label: "Repository", Name: "user-service", Properties: map[string]interface{}{"name": "user-service", "url": "https://github.com/acme/user-service", "language": "Go", "visibility": "private"}},
		{ID: "repo-frontend", Type: "Repository", Label: "Repository", Name: "frontend", Properties: map[string]interface{}{"name": "frontend", "url": "https://github.com/orbit/frontend", "language": "TypeScript", "visibility": "public"}},
		{ID: "repo-inventory-service", Type: "Repository", Label: "Repository", Name: "inventory-service", Properties: map[string]interface{}{"name": "inventory-service", "url": "https://github.com/acme/inventory-service", "language": "Python", "visibility": "private"}},
		{ID: "repo-analytics-service", Type: "Repository", Label: "Repository", Name: "analytics-service", Properties: map[string]interface{}{"name": "analytics-service", "url": "https://github.com/orbit/analytics-service", "language": "Python", "visibility": "private"}},
		{ID: "repo-order-service", Type: "Repository", Label: "Repository", Name: "order-service", Properties: map[string]interface{}{"name": "order-service", "url": "https://github.com/acme/order-service", "language": "Go", "visibility": "private"}},
		{ID: "repo-backend-api", Type: "Repository", Label: "Repository", Name: "backend-api", Properties: map[string]interface{}{"name": "backend-api", "url": "https://github.com/acme/backend-api", "language": "Go", "visibility": "private"}},
		{ID: "repo-notification-api", Type: "Repository", Label: "Repository", Name: "notification-api", Properties: map[string]interface{}{"name": "notification-api", "url": "https://github.com/nova/notification-api", "language": "Go", "visibility": "private"}},
		{ID: "repo-auth-service", Type: "Repository", Label: "Repository", Name: "auth-service", Properties: map[string]interface{}{"name": "auth-service", "url": "https://github.com/acme/auth-service", "language": "Go", "visibility": "private"}},
		{ID: "repo-search-service", Type: "Repository", Label: "Repository", Name: "search-service", Properties: map[string]interface{}{"name": "search-service", "url": "https://github.com/orbit/search-service", "language": "Python", "visibility": "private"}},
		{ID: "repo-cart-service", Type: "Repository", Label: "Repository", Name: "cart-service", Properties: map[string]interface{}{"name": "cart-service", "url": "https://github.com/acme/cart-service", "language": "Go", "visibility": "private"}},
		{ID: "repo-shipping-service", Type: "Repository", Label: "Repository", Name: "shipping-service", Properties: map[string]interface{}{"name": "shipping-service", "url": "https://github.com/nova/shipping-service", "language": "Java", "visibility": "private"}},
		{ID: "repo-recommendation-service", Type: "Repository", Label: "Repository", Name: "recommendation-service", Properties: map[string]interface{}{"name": "recommendation-service", "url": "https://github.com/orbit/recommendation-service", "language": "Python", "visibility": "private"}},
		{ID: "repo-billing-service", Type: "Repository", Label: "Repository", Name: "billing-service", Properties: map[string]interface{}{"name": "billing-service", "url": "https://github.com/nova/billing-service", "language": "Go", "visibility": "private"}},
		{ID: "repo-catalog-service", Type: "Repository", Label: "Repository", Name: "catalog-service", Properties: map[string]interface{}{"name": "catalog-service", "url": "https://github.com/acme/catalog-service", "language": "Go", "visibility": "private"}},
		{ID: "repo-gateway", Type: "Repository", Label: "Repository", Name: "api-gateway", Properties: map[string]interface{}{"name": "api-gateway", "url": "https://github.com/acme/api-gateway", "language": "Go", "visibility": "private"}},
		{ID: "repo-worker", Type: "Repository", Label: "Repository", Name: "background-worker", Properties: map[string]interface{}{"name": "background-worker", "url": "https://github.com/orbit/background-worker", "language": "Python", "visibility": "private"}},
		{ID: "repo-mobile-app", Type: "Repository", Label: "Repository", Name: "mobile-app", Properties: map[string]interface{}{"name": "mobile-app", "url": "https://github.com/orbit/mobile-app", "language": "Swift", "visibility": "public"}},
		{ID: "repo-admin-panel", Type: "Repository", Label: "Repository", Name: "admin-panel", Properties: map[string]interface{}{"name": "admin-panel", "url": "https://github.com/acme/admin-panel", "language": "TypeScript", "visibility": "private"}},
		{ID: "repo-event-bus", Type: "Repository", Label: "Repository", Name: "event-bus", Properties: map[string]interface{}{"name": "event-bus", "url": "https://github.com/acme/event-bus", "language": "Go", "visibility": "private"}},
		{ID: "repo-file-storage", Type: "Repository", Label: "Repository", Name: "file-storage", Properties: map[string]interface{}{"name": "file-storage", "url": "https://github.com/nova/file-storage", "language": "Go", "visibility": "private"}},
		{ID: "repo-cache-layer", Type: "Repository", Label: "Repository", Name: "cache-layer", Properties: map[string]interface{}{"name": "cache-layer", "url": "https://github.com/acme/cache-layer", "language": "Go", "visibility": "private"}},
		{ID: "repo-monitoring", Type: "Repository", Label: "Repository", Name: "monitoring", Properties: map[string]interface{}{"name": "monitoring", "url": "https://github.com/acme/monitoring", "language": "Go", "visibility": "private"}},
		{ID: "repo-logger", Type: "Repository", Label: "Repository", Name: "logger-lib", Properties: map[string]interface{}{"name": "logger-lib", "url": "https://github.com/acme/logger-lib", "language": "Go", "visibility": "public"}},
		{ID: "repo-sdk", Type: "Repository", Label: "Repository", Name: "client-sdk", Properties: map[string]interface{}{"name": "client-sdk", "url": "https://github.com/orbit/client-sdk", "language": "TypeScript", "visibility": "public"}},
		{ID: "repo-webhooks", Type: "Repository", Label: "Repository", Name: "webhook-handler", Properties: map[string]interface{}{"name": "webhook-handler", "url": "https://github.com/nova/webhook-handler", "language": "Go", "visibility": "private"}},
		{ID: "repo-queue-worker", Type: "Repository", Label: "Repository", Name: "queue-worker", Properties: map[string]interface{}{"name": "queue-worker", "url": "https://github.com/orbit/queue-worker", "language": "Go", "visibility": "private"}},
		{ID: "repo-feature-flags", Type: "Repository", Label: "Repository", Name: "feature-flags", Properties: map[string]interface{}{"name": "feature-flags", "url": "https://github.com/acme/feature-flags", "language": "Go", "visibility": "private"}},
	}
	for _, r := range repos {
		s.addNode(r)
	}

	packages := []models.Node{
		{ID: "pkg-postgres", Type: "Package", Label: "Package", Name: "PostgreSQL", Properties: map[string]interface{}{"name": "PostgreSQL", "version": "16.2", "ecosystem": "apt"}},
		{ID: "pkg-redis", Type: "Package", Label: "Package", Name: "Redis", Properties: map[string]interface{}{"name": "Redis", "version": "7.2", "ecosystem": "apt"}},
		{ID: "pkg-kafka", Type: "Package", Label: "Package", Name: "kafka-client", Properties: map[string]interface{}{"name": "kafka-client", "version": "3.6", "ecosystem": "npm"}},
		{ID: "pkg-gin", Type: "Package", Label: "Package", Name: "gin", Properties: map[string]interface{}{"name": "gin", "version": "1.9.1", "ecosystem": "go"}},
		{ID: "pkg-grpc", Type: "Package", Label: "Package", Name: "grpc", Properties: map[string]interface{}{"name": "grpc", "version": "1.60", "ecosystem": "go"}},
		{ID: "pkg-stripe", Type: "Package", Label: "Package", Name: "stripe-sdk", Properties: map[string]interface{}{"name": "stripe-sdk", "version": "74.0", "ecosystem": "go"}},
		{ID: "pkg-react", Type: "Package", Label: "Package", Name: "React", Properties: map[string]interface{}{"name": "React", "version": "18.2", "ecosystem": "npm"}},
		{ID: "pkg-axios", Type: "Package", Label: "Package", Name: "axios", Properties: map[string]interface{}{"name": "axios", "version": "1.6", "ecosystem": "npm"}},
		{ID: "pkg-httpclient", Type: "Package", Label: "Package", Name: "http-client", Properties: map[string]interface{}{"name": "http-client", "version": "2.1", "ecosystem": "go"}},
		{ID: "pkg-pgx", Type: "Package", Label: "Package", Name: "pgx", Properties: map[string]interface{}{"name": "pgx", "version": "5.5", "ecosystem": "go"}},
		{ID: "pkg-gorm", Type: "Package", Label: "Package", Name: "gorm", Properties: map[string]interface{}{"name": "gorm", "version": "2.4", "ecosystem": "go"}},
		{ID: "pkg-jwt", Type: "Package", Label: "Package", Name: "jwt-go", Properties: map[string]interface{}{"name": "jwt-go", "version": "4.5", "ecosystem": "go"}},
		{ID: "pkg-elasticsearch", Type: "Package", Label: "Package", Name: "elasticsearch", Properties: map[string]interface{}{"name": "elasticsearch", "version": "8.11", "ecosystem": "go"}},
		{ID: "pkg-prometheus", Type: "Package", Label: "Package", Name: "prometheus-client", Properties: map[string]interface{}{"name": "prometheus-client", "version": "1.19", "ecosystem": "go"}},
		{ID: "pkg-zap", Type: "Package", Label: "Package", Name: "zap", Properties: map[string]interface{}{"name": "zap", "version": "1.26", "ecosystem": "go"}},
		{ID: "pkg-uuid", Type: "Package", Label: "Package", Name: "uuid", Properties: map[string]interface{}{"name": "uuid", "version": "1.5", "ecosystem": "go"}},
		{ID: "pkg-lodash", Type: "Package", Label: "Package", Name: "lodash", Properties: map[string]interface{}{"name": "lodash", "version": "4.17", "ecosystem": "npm"}},
		{ID: "pkg-express", Type: "Package", Label: "Package", Name: "express", Properties: map[string]interface{}{"name": "express", "version": "4.18", "ecosystem": "npm"}},
		{ID: "pkg-typescript", Type: "Package", Label: "Package", Name: "typescript", Properties: map[string]interface{}{"name": "typescript", "version": "5.3", "ecosystem": "npm"}},
		{ID: "pkg-vite", Type: "Package", Label: "Package", Name: "vite", Properties: map[string]interface{}{"name": "vite", "version": "5.0", "ecosystem": "npm"}},
		{ID: "pkg-tailwind", Type: "Package", Label: "Package", Name: "tailwindcss", Properties: map[string]interface{}{"name": "tailwindcss", "version": "3.4", "ecosystem": "npm"}},
		{ID: "pkg-nextjs", Type: "Package", Label: "Package", Name: "next", Properties: map[string]interface{}{"name": "next", "version": "14.1", "ecosystem": "npm"}},
		{ID: "pkg-prisma", Type: "Package", Label: "Package", Name: "prisma", Properties: map[string]interface{}{"name": "prisma", "version": "5.8", "ecosystem": "npm"}},
		{ID: "pkg-zod", Type: "Package", Label: "Package", Name: "zod", Properties: map[string]interface{}{"name": "zod", "version": "3.22", "ecosystem": "npm"}},
		{ID: "pkg-neo4j-driver", Type: "Package", Label: "Package", Name: "neo4j-driver", Properties: map[string]interface{}{"name": "neo4j-driver", "version": "5.15", "ecosystem": "go"}},
		{ID: "pkg-cors", Type: "Package", Label: "Package", Name: "cors", Properties: map[string]interface{}{"name": "cors", "version": "2.8", "ecosystem": "go"}},
		{ID: "pkg-testify", Type: "Package", Label: "Package", Name: "testify", Properties: map[string]interface{}{"name": "testify", "version": "1.8", "ecosystem": "go"}},
		{ID: "pkg-docker", Type: "Package", Label: "Package", Name: "docker-client", Properties: map[string]interface{}{"name": "docker-client", "version": "24.0", "ecosystem": "go"}},
		{ID: "pkg-kubernetes", Type: "Package", Label: "Package", Name: "k8s-client", Properties: map[string]interface{}{"name": "k8s-client", "version": "0.29", "ecosystem": "go"}},
		{ID: "pkg-opentelemetry", Type: "Package", Label: "Package", Name: "opentelemetry", Properties: map[string]interface{}{"name": "opentelemetry", "version": "1.21", "ecosystem": "go"}},
		{ID: "pkg-grpc-gateway", Type: "Package", Label: "Package", Name: "grpc-gateway", Properties: map[string]interface{}{"name": "grpc-gateway", "version": "2.18", "ecosystem": "go"}},
		{ID: "pkg-protobuf", Type: "Package", Label: "Package", Name: "protobuf", Properties: map[string]interface{}{"name": "protobuf", "version": "1.32", "ecosystem": "go"}},
		{ID: "pkg-viper", Type: "Package", Label: "Package", Name: "viper", Properties: map[string]interface{}{"name": "viper", "version": "1.18", "ecosystem": "go"}},
		{ID: "pkg-cobra", Type: "Package", Label: "Package", Name: "cobra", Properties: map[string]interface{}{"name": "cobra", "version": "1.8", "ecosystem": "go"}},
		{ID: "pkg-echo", Type: "Package", Label: "Package", Name: "echo", Properties: map[string]interface{}{"name": "echo", "version": "4.11", "ecosystem": "go"}},
		{ID: "pkg-chi", Type: "Package", Label: "Package", Name: "chi", Properties: map[string]interface{}{"name": "chi", "version": "5.0", "ecosystem": "go"}},
		{ID: "pkg-wire", Type: "Package", Label: "Package", Name: "wire", Properties: map[string]interface{}{"name": "wire", "version": "0.6", "ecosystem": "go"}},
		{ID: "pkg-fiber", Type: "Package", Label: "Package", Name: "fiber", Properties: map[string]interface{}{"name": "fiber", "version": "2.52", "ecosystem": "go"}},
		{ID: "pkg-bcrypt", Type: "Package", Label: "Package", Name: "bcrypt", Properties: map[string]interface{}{"name": "bcrypt", "version": "0.15", "ecosystem": "go"}},
		{ID: "pkg-nodepg", Type: "Package", Label: "Package", Name: "node-postgres", Properties: map[string]interface{}{"name": "node-postgres", "version": "8.11", "ecosystem": "npm"}},
	}
	for _, p := range packages {
		s.addNode(p)
	}

	services := []models.Node{
		{ID: "svc-checkout", Type: "Service", Label: "Service", Name: "Checkout Service", Properties: map[string]interface{}{"name": "Checkout Service", "environment": "production", "description": "Handles checkout flow"}},
		{ID: "svc-payments", Type: "Service", Label: "Service", Name: "Payment Service", Properties: map[string]interface{}{"name": "Payment Service", "environment": "production", "description": "Processes payments via Stripe"}},
		{ID: "svc-notifications", Type: "Service", Label: "Service", Name: "Notification Service", Properties: map[string]interface{}{"name": "Notification Service", "environment": "production", "description": "Sends emails/push notifications"}},
		{ID: "svc-orders", Type: "Service", Label: "Service", Name: "Order Service", Properties: map[string]interface{}{"name": "Order Service", "environment": "production", "description": "Manages orders"}},
		{ID: "svc-users", Type: "Service", Label: "Service", Name: "User Service", Properties: map[string]interface{}{"name": "User Service", "environment": "production", "description": "User accounts and auth"}},
		{ID: "svc-inventory", Type: "Service", Label: "Service", Name: "Inventory Service", Properties: map[string]interface{}{"name": "Inventory Service", "environment": "production", "description": "Stock management"}},
		{ID: "svc-analytics", Type: "Service", Label: "Service", Name: "Analytics Service", Properties: map[string]interface{}{"name": "Analytics Service", "environment": "production", "description": "Event analytics"}},
		{ID: "svc-backend-api", Type: "Service", Label: "Service", Name: "Backend API", Properties: map[string]interface{}{"name": "Backend API", "environment": "production", "description": "Main API gateway"}},
		{ID: "svc-auth", Type: "Service", Label: "Service", Name: "Auth Service", Properties: map[string]interface{}{"name": "Auth Service", "environment": "production", "description": "Authentication"}},
		{ID: "svc-search", Type: "Service", Label: "Service", Name: "Search Service", Properties: map[string]interface{}{"name": "Search Service", "environment": "production", "description": "Product search via Elasticsearch"}},
		{ID: "svc-cart", Type: "Service", Label: "Service", Name: "Cart Service", Properties: map[string]interface{}{"name": "Cart Service", "environment": "production", "description": "Shopping cart"}},
		{ID: "svc-shipping", Type: "Service", Label: "Service", Name: "Shipping Service", Properties: map[string]interface{}{"name": "Shipping Service", "environment": "production", "description": "Shipping logistics"}},
		{ID: "svc-catalog", Type: "Service", Label: "Service", Name: "Catalog Service", Properties: map[string]interface{}{"name": "Catalog Service", "environment": "production", "description": "Product catalog"}},
		{ID: "svc-gateway", Type: "Service", Label: "Service", Name: "API Gateway", Properties: map[string]interface{}{"name": "API Gateway", "environment": "production", "description": "Edge gateway"}},
		{ID: "svc-worker", Type: "Service", Label: "Service", Name: "Background Worker", Properties: map[string]interface{}{"name": "Background Worker", "environment": "production", "description": "Async jobs"}},
	}
	for _, svc := range services {
		s.addNode(svc)
	}

	// Relationships — carefully crafted for demo paths
	edge := func(src, typ, dst string, props map[string]interface{}) {
		s.addEdge(models.Edge{ID: fmt.Sprintf("%s-%s-%s", src, typ, dst), Source: src, Target: dst, Type: typ, Properties: props})
	}

	// Repository -> Package DEPENDS_ON
	repoDeps := [][3]string{
		{"repo-checkout-api", "pkg-postgres", "16.2"}, {"repo-checkout-api", "pkg-gin", "1.9.1"}, {"repo-checkout-api", "pkg-pgx", "5.5"}, {"repo-checkout-api", "pkg-stripe", "74.0"},
		{"repo-payment-service", "pkg-postgres", "16.2"}, {"repo-payment-service", "pkg-stripe", "74.0"}, {"repo-payment-service", "pkg-grpc", "1.60"},
		{"repo-notification-service", "pkg-redis", "7.2"}, {"repo-notification-service", "pkg-kafka", "3.6"}, {"repo-notification-service", "pkg-gin", "1.9.1"},
		{"repo-backend-api", "pkg-postgres", "16.2"}, {"repo-backend-api", "pkg-redis", "7.2"}, {"repo-backend-api", "pkg-gin", "1.9.1"}, {"repo-backend-api", "pkg-jwt", "4.5"},
		{"repo-order-service", "pkg-postgres", "16.2"}, {"repo-order-service", "pkg-kafka", "3.6"}, {"repo-order-service", "pkg-grpc", "1.60"},
		{"repo-notification-api", "pkg-redis", "7.2"}, {"repo-notification-api", "pkg-kafka", "3.6"},
		{"repo-user-service", "pkg-postgres", "16.2"}, {"repo-user-service", "pkg-jwt", "4.5"}, {"repo-user-service", "pkg-bcrypt", "0.15"},
		{"repo-frontend", "pkg-react", "18.2"}, {"repo-frontend", "pkg-axios", "1.6"}, {"repo-frontend", "pkg-nextjs", "14.1"}, {"repo-frontend", "pkg-tailwind", "3.4"},
		{"repo-analytics-service", "pkg-kafka", "3.6"}, {"repo-analytics-service", "pkg-elasticsearch", "8.11"},
		{"repo-inventory-service", "pkg-postgres", "16.2"}, {"repo-inventory-service", "pkg-redis", "7.2"},
		{"repo-auth-service", "pkg-postgres", "16.2"}, {"repo-auth-service", "pkg-jwt", "4.5"},
		{"repo-search-service", "pkg-elasticsearch", "8.11"}, {"repo-search-service", "pkg-redis", "7.2"},
		{"repo-cart-service", "pkg-redis", "7.2"}, {"repo-cart-service", "pkg-postgres", "16.2"},
		{"repo-shipping-service", "pkg-kafka", "3.6"}, {"repo-shipping-service", "pkg-grpc", "1.60"},
		{"repo-catalog-service", "pkg-postgres", "16.2"}, {"repo-catalog-service", "pkg-elasticsearch", "8.11"},
		{"repo-gateway", "pkg-grpc", "1.60"}, {"repo-gateway", "pkg-grpc-gateway", "2.18"}, {"repo-gateway", "pkg-neo4j-driver", "5.15"},
		{"repo-event-bus", "pkg-kafka", "3.6"}, {"repo-event-bus", "pkg-protobuf", "1.32"},
		{"repo-cache-layer", "pkg-redis", "7.2"}, {"repo-cache-layer", "pkg-gin", "1.9.1"},
		{"repo-monitoring", "pkg-prometheus", "1.19"}, {"repo-monitoring", "pkg-zap", "1.26"}, {"repo-monitoring", "pkg-opentelemetry", "1.21"},
		{"repo-admin-panel", "pkg-react", "18.2"}, {"repo-admin-panel", "pkg-axios", "1.6"},
		{"repo-mobile-app", "pkg-axios", "1.6"}, {"repo-mobile-app", "pkg-react", "18.2"},
		{"repo-billing-service", "pkg-postgres", "16.2"}, {"repo-billing-service", "pkg-stripe", "74.0"},
		{"repo-queue-worker", "pkg-kafka", "3.6"}, {"repo-queue-worker", "pkg-redis", "7.2"},
	}
	for _, d := range repoDeps {
		edge(d[0], "DEPENDS_ON", d[1], map[string]interface{}{"version": d[2], "direct": true})
	}

	// Package -> Package
	pkgDeps := [][3]string{
		{"pkg-stripe", "pkg-httpclient", "2.1"},
		{"pkg-gin", "pkg-httpclient", "2.1"},
		{"pkg-pgx", "pkg-postgres", "16.2"},
		{"pkg-gorm", "pkg-postgres", "16.2"},
		{"pkg-grpc-gateway", "pkg-grpc", "1.60"},
		{"pkg-grpc-gateway", "pkg-protobuf", "1.32"},
		{"pkg-nextjs", "pkg-react", "18.2"},
		{"pkg-prisma", "pkg-postgres", "16.2"},
		{"pkg-elasticsearch", "pkg-httpclient", "2.1"},
	}
	for _, d := range pkgDeps {
		edge(d[0], "DEPENDS_ON", d[1], map[string]interface{}{"version": d[2]})
	}

	// Repository -> Service DEPLOYED_AS
	deployed := [][2]string{
		{"repo-checkout-api", "svc-checkout"}, {"repo-payment-service", "svc-payments"}, {"repo-notification-service", "svc-notifications"},
		{"repo-order-service", "svc-orders"}, {"repo-user-service", "svc-users"}, {"repo-inventory-service", "svc-inventory"},
		{"repo-analytics-service", "svc-analytics"}, {"repo-backend-api", "svc-backend-api"}, {"repo-auth-service", "svc-auth"},
		{"repo-search-service", "svc-search"}, {"repo-cart-service", "svc-cart"}, {"repo-shipping-service", "svc-shipping"},
		{"repo-catalog-service", "svc-catalog"}, {"repo-gateway", "svc-gateway"}, {"repo-queue-worker", "svc-worker"},
		{"repo-notification-api", "svc-notifications"}, {"repo-billing-service", "svc-payments"},
	}
	for _, d := range deployed {
		edge(d[0], "DEPLOYED_AS", d[1], nil)
	}

	// Service -> Service CALLS (creates the impact chains)
	calls := [][2]string{
		{"svc-checkout", "svc-payments"}, {"svc-checkout", "svc-orders"}, {"svc-checkout", "svc-inventory"},
		{"svc-orders", "svc-notifications"}, {"svc-orders", "svc-payments"}, {"svc-orders", "svc-catalog"},
		{"svc-backend-api", "svc-orders"}, {"svc-backend-api", "svc-users"}, {"svc-backend-api", "svc-cart"},
		{"svc-cart", "svc-checkout"}, {"svc-gateway", "svc-backend-api"}, {"svc-gateway", "svc-auth"},
		{"svc-notifications", "svc-analytics"}, {"svc-payments", "svc-notifications"},
		{"svc-catalog", "svc-search"}, {"svc-search", "svc-analytics"}, {"svc-worker", "svc-notifications"},
		{"svc-shipping", "svc-orders"}, {"svc-users", "svc-auth"},
	}
	for _, c := range calls {
		edge(c[0], "CALLS", c[1], nil)
	}

	// Team OWNS Repository/Service
	owns := [][2]string{
		{"team-checkout", "repo-checkout-api"}, {"team-checkout", "svc-checkout"},
		{"team-payments", "repo-payment-service"}, {"team-payments", "svc-payments"}, {"team-payments", "repo-billing-service"},
		{"team-notifications", "repo-notification-service"}, {"team-notifications", "repo-notification-api"}, {"team-notifications", "svc-notifications"},
		{"team-platform", "repo-backend-api"}, {"team-platform", "svc-backend-api"}, {"team-platform", "repo-gateway"}, {"team-platform", "svc-gateway"},
		{"team-platform", "repo-order-service"}, {"team-platform", "svc-orders"},
		{"team-infra", "repo-cache-layer"}, {"team-infra", "repo-monitoring"}, {"team-infra", "repo-event-bus"},
		{"team-frontend", "repo-frontend"}, {"team-frontend", "repo-mobile-app"}, {"team-frontend", "repo-admin-panel"},
		{"team-data", "repo-analytics-service"}, {"team-data", "svc-analytics"}, {"team-data", "repo-search-service"}, {"team-data", "svc-search"},
		{"team-security", "repo-auth-service"}, {"team-security", "svc-auth"},
		{"team-checkout", "repo-cart-service"}, {"team-checkout", "svc-cart"}, {"team-checkout", "repo-inventory-service"}, {"team-checkout", "svc-inventory"},
		{"team-platform", "repo-catalog-service"}, {"team-platform", "svc-catalog"}, {"team-platform", "repo-shipping-service"}, {"team-platform", "svc-shipping"},
		{"team-data", "repo-queue-worker"}, {"team-data", "svc-worker"},
	}
	for _, o := range owns {
		edge(o[0], "OWNS", o[1], nil)
	}

	// Repository BELONGS_TO Organization
	belongs := [][2]string{
		{"repo-checkout-api", "org-acme"}, {"repo-backend-api", "org-acme"}, {"repo-user-service", "org-acme"}, {"repo-cart-service", "org-acme"}, {"repo-catalog-service", "org-acme"}, {"repo-gateway", "org-acme"}, {"repo-auth-service", "org-acme"}, {"repo-cache-layer", "org-acme"}, {"repo-monitoring", "org-acme"}, {"repo-admin-panel", "org-acme"}, {"repo-feature-flags", "org-acme"}, {"repo-logger", "org-acme"},
		{"repo-payment-service", "org-nova"}, {"repo-notification-service", "org-nova"}, {"repo-notification-api", "org-nova"}, {"repo-order-service", "org-acme"}, {"repo-shipping-service", "org-nova"}, {"repo-billing-service", "org-nova"}, {"repo-file-storage", "org-nova"}, {"repo-webhooks", "org-nova"},
		{"repo-frontend", "org-orbit"}, {"repo-analytics-service", "org-orbit"}, {"repo-inventory-service", "org-acme"}, {"repo-search-service", "org-orbit"}, {"repo-event-bus", "org-acme"}, {"repo-worker", "org-orbit"}, {"repo-mobile-app", "org-orbit"}, {"repo-sdk", "org-orbit"}, {"repo-queue-worker", "org-orbit"}, {"repo-recommendation-service", "org-orbit"},
	}
	for _, b := range belongs {
		edge(b[0], "BELONGS_TO", b[1], nil)
	}

	// Developer CONTRIBUTED_TO / MAINTAINS
	contribs := [][2]string{
		{"dev-1", "repo-checkout-api"}, {"dev-2", "repo-checkout-api"}, {"dev-6", "repo-backend-api"}, {"dev-10", "repo-backend-api"},
		{"dev-3", "repo-payment-service"}, {"dev-4", "repo-notification-service"}, {"dev-7", "repo-user-service"}, {"dev-5", "repo-frontend"},
		{"dev-8", "repo-order-service"}, {"dev-11", "repo-notification-api"}, {"dev-12", "repo-auth-service"}, {"dev-13", "repo-analytics-service"},
		{"dev-14", "repo-inventory-service"}, {"dev-15", "repo-search-service"}, {"dev-16", "repo-cart-service"}, {"dev-17", "repo-catalog-service"},
		{"dev-18", "repo-gateway"}, {"dev-19", "repo-billing-service"}, {"dev-20", "repo-queue-worker"},
	}
	for _, c := range contribs {
		edge(c[0], "CONTRIBUTED_TO", c[1], map[string]interface{}{"commits": 42, "last_contribution": "2024-03-15"})
	}
	maintains := [][2]string{
		{"dev-1", "repo-checkout-api"}, {"dev-3", "repo-payment-service"}, {"dev-4", "repo-notification-service"}, {"dev-6", "repo-backend-api"}, {"dev-5", "repo-frontend"},
	}
	for _, m := range maintains {
		edge(m[0], "MAINTAINS", m[1], nil)
	}

	// recent edges for dashboard
	s.recent = []models.RecentRel{
		{Source: "checkout-api", Target: "PostgreSQL", Type: "DEPENDS_ON", Detail: "v16.2"},
		{Source: "payment-service", Target: "stripe-sdk", Type: "DEPENDS_ON", Detail: "v74.0"},
		{Source: "backend-api", Target: "Redis", Type: "DEPENDS_ON", Detail: "v7.2"},
		{Source: "order-service", Target: "kafka-client", Type: "DEPENDS_ON", Detail: "v3.6"},
		{Source: "frontend", Target: "React", Type: "DEPENDS_ON", Detail: "v18.2"},
		{Source: "analytics-service", Target: "elasticsearch", Type: "DEPENDS_ON", Detail: "v8.11"},
		{Source: "checkout-api", Target: "Checkout Service", Type: "DEPLOYED_AS", Detail: ""},
		{Source: "Checkout Service", Target: "Payment Service", Type: "CALLS", Detail: ""},
		{Source: "Order Service", Target: "Notification Service", Type: "CALLS", Detail: ""},
		{Source: "Backend API", Target: "Order Service", Type: "CALLS", Detail: ""},
	}
}
