// check.js — Custom Playwright browser check script.
//
// This script runs inside Oack's sandboxed browser checker. You get two
// globals: `page` (a Playwright Page) and `context` (a BrowserContext).
// No require/import, no filesystem, no process access — by design.
//
// The checker evaluates pass/fail based on:
//   1. No uncaught exceptions during execution.
//   2. Final page status code is in allowed_status_codes (default: 2xx, 3xx).
//   3. Total execution time < timeout_ms.
//   4. Console error count < console_error_threshold.
//   5. Resource error count < resource_error_threshold.
//
// Environment variables from `script_env` are available as globals.
// For example, if you set { key: "LOGIN_EMAIL", value: "test@example.com" },
// you can reference it as `LOGIN_EMAIL` in this script.

// ---- Example: Login flow + dashboard health check ----

// Step 1: Navigate to the login page.
await page.goto("https://app.example.com/login");

// Step 2: Fill in the login form.
await page.fill('input[name="email"]', LOGIN_EMAIL);
await page.fill('input[name="password"]', LOGIN_PASSWORD);

// Step 3: Submit and wait for navigation.
await page.click('button[type="submit"]');
await page.waitForURL("**/dashboard", { timeout: 10000 });

// Step 4: Verify the dashboard loaded correctly.
const heading = await page.textContent("h1");
if (!heading || !heading.includes("Dashboard")) {
  throw new Error(`Expected "Dashboard" heading, got: "${heading}"`);
}

// Step 5: Check that a critical API widget rendered.
await page.waitForSelector('[data-testid="revenue-widget"]', { timeout: 5000 });

// Step 6: Take a screenshot (optional — the checker captures one automatically,
// but you can take additional ones at specific points).
// await page.screenshot({ path: '/tmp/dashboard.png' });

// If the script completes without throwing, the check passes.
