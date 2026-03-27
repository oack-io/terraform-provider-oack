// check.js — Custom Playwright browser check script.
//
// Export a single async function that receives:
//   - page:    Playwright Page object
//   - context: Object with env vars from script_env as properties
//              (e.g. context.LOGIN_EMAIL, context.LOGIN_PASSWORD)
//
// The function runs inside Oack's sandboxed browser checker.
// No require/import, no filesystem, no process access — by design.
//
// Pass/fail is determined by:
//   1. The function completes without throwing.
//   2. Final page status code is in allowed_status_codes (default: 2xx, 3xx).
//   3. Total execution time < timeout_ms.
//   4. Console error count < console_error_threshold.
//   5. Resource error count < resource_error_threshold.

module.exports = async function (page, context) {
  // ---- Example: Login flow + dashboard health check ----

  // Step 1: Navigate to the login page.
  await page.goto("https://app.example.com/login");

  // Step 2: Fill in the login form using env vars from script_env.
  await page.fill('input[name="email"]', context.LOGIN_EMAIL);
  await page.fill('input[name="password"]', context.LOGIN_PASSWORD);

  // Step 3: Submit and wait for navigation.
  await page.click('button[type="submit"]');
  await page.waitForURL("**/dashboard", { timeout: 10000 });

  // Step 4: Verify the dashboard loaded correctly.
  const heading = await page.textContent("h1");
  if (!heading || !heading.includes("Dashboard")) {
    throw new Error(`Expected "Dashboard" heading, got: "${heading}"`);
  }

  // Step 5: Check that a critical API widget rendered.
  await page.waitForSelector('[data-testid="revenue-widget"]', {
    timeout: 5000,
  });

  // If the function completes without throwing, the check passes.
};
