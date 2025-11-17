vX.Y.0

Below is a list of all new features and changes included in this minor release.

Summary
<One or two concise lines: highlight primary new capabilities + scope.>

New Features
Backend
* <New endpoint: GET /... (purpose)>  
* <New domain model / entity: X (fields: ...)>  
* <Background worker / scheduled task for Y>

Frontend
* <New page/component: ProductReviewList with filtering>
* <New route /profile for managing public identity>
* <State management improvement: moved to context/store>

Integrations (optional)
* <Added payment gateway integration (ProviderName)>  
* <Webhook endpoint for external system>

Improvements
Backend
* <Refactored service layering for clarity>
* <Optimized query for listing products (N+1 removed)>
* <Stronger validation around field X>

Frontend
* <Improved accessibility: ARIA labels added>
* <Better loading states for product grid>

Fixes
* <Short bullet list of notable bug fixes>
* <Backend: fix race in cache invalidation>
* <Frontend: fix layout shift on mobile>

Performance (optional)
* <Reduced average response time for /products by ~30%>
* <Cut bundle size from 450KB -> 320KB>

Security (optional)
* <Upgraded dependency to patch CVE-XXXX-YYYY>
* <Added rate limiting to auth endpoints>

Deprecations (optional)
* <Endpoint /v1/legacy marked deprecated;>

Migrations (optional)
* None
* <2025MMDDHHMMSS_add_field_to_table.up.sql>
* <2025MMDDHHMMSS_create_new_table.up.sql>

Upgrade Notes (optional)
* Run pending database migrations.
* Clear application cache.
* Rebuild frontend assets.

Important commits
* <hash> feat(...): add main feature A
* <hash> feat(...): introduce endpoint B
* <hash> refactor(...): reorganize modules
* <hash> perf(...): optimize query C
* <hash> fix(...): resolve issue D

Notes (optional)
* <UI field now has maxLength=500>
* <Deprecation warnings added for legacy endpoint>

Dependencies (optional)
* <go module foo v1.2.3 -> v1.3.0>
* <npm package bar 2.5.0 -> 2.6.0>

Tooling (optional)
* <Added lint rule set X>
* <Introduced pre-commit hooks>