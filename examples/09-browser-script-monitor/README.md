# Example 09: Browser Monitor with Custom Playwright Script

This example shows how to write a custom Playwright script, test it locally,
and then deploy it as an Oack browser monitor using Terraform.

## The Script

[`check.js`](check.js) implements a login → dashboard health check:

1. Navigate to the login page
2. Fill email and password
3. Submit the form and wait for redirect
4. Assert the dashboard heading
5. Assert a critical widget rendered

## Run Locally with Node.js

You can develop and debug the script on your machine before deploying it.

### Prerequisites

```bash
# Install Node.js (18+) and Playwright
npm init -y
npm install playwright
npx playwright install chromium
```

### Create a local test harness

Create `test-local.js` — a thin wrapper that sets up the same globals your
script will have in the Oack sandbox:

```js
// test-local.js — run check.js locally with Playwright
const { chromium } = require("playwright");

(async () => {
  // Set environment variables (same as script_env in Terraform).
  const LOGIN_EMAIL = process.env.LOGIN_EMAIL || "test@example.com";
  const LOGIN_PASSWORD = process.env.LOGIN_PASSWORD || "s3cret";

  const browser = await chromium.launch({ headless: true });
  const context = await browser.newContext({
    viewport: { width: 1920, height: 1080 },
  });
  const page = await context.newPage();

  try {
    // Make env vars available as globals (matches Oack sandbox behavior).
    await page.evaluate(
      ([email, password]) => {
        globalThis.LOGIN_EMAIL = email;
        globalThis.LOGIN_PASSWORD = password;
      },
      [LOGIN_EMAIL, LOGIN_PASSWORD]
    );

    // Run the check script in the page context.
    // Note: in production, Oack runs the script in a Node.js sandbox with
    // `page` and `context` as globals. Locally, we just inline the steps.

    // --- Paste your check.js steps here, or use eval: ---
    await page.goto("https://app.example.com/login");
    await page.fill('input[name="email"]', LOGIN_EMAIL);
    await page.fill('input[name="password"]', LOGIN_PASSWORD);
    await page.click('button[type="submit"]');
    await page.waitForURL("**/dashboard", { timeout: 10000 });

    const heading = await page.textContent("h1");
    if (!heading || !heading.includes("Dashboard")) {
      throw new Error(`Expected "Dashboard" heading, got: "${heading}"`);
    }

    await page.waitForSelector('[data-testid="revenue-widget"]', {
      timeout: 5000,
    });

    console.log("PASS: Login flow completed successfully");
  } catch (err) {
    console.error("FAIL:", err.message);
    await page.screenshot({ path: "failure-screenshot.png" });
    process.exit(1);
  } finally {
    await browser.close();
  }
})();
```

### Run it

```bash
LOGIN_EMAIL="test@example.com" LOGIN_PASSWORD="s3cret" node test-local.js
```

You'll see either `PASS: Login flow completed successfully` or a failure
message with a screenshot saved to `failure-screenshot.png`.

### Tips for local development

- **Run headed** for debugging: change `headless: true` to `headless: false`
- **Slow down** for visibility: add `slowMo: 100` to `chromium.launch()`
- **Record**: use `npx playwright codegen https://app.example.com/login` to
  generate the selector/action code automatically
- **Timeout**: the default Oack timeout is 30s — keep your script well within that

## Deploy to Oack with Terraform

Once the script works locally, deploy it:

```bash
export OACK_API_KEY="oack_acc_xxxxxxxxxxxxxxxxxxxx"
export OACK_ACCOUNT_ID="your-account-uuid"
export TF_VAR_login_email="test@example.com"
export TF_VAR_login_password="s3cret"

terraform init
terraform plan
terraform apply
```

Oack will now run `check.js` every 5 minutes from its browser checker nodes.
You get:
- **Web Vitals** (LCP, FCP, CLS, TTFB) for every run
- **Screenshots** captured automatically
- **HAR files** with full resource waterfall
- **Console messages** (errors and warnings)
- **Alerts** via Slack (or any configured channel) when the flow breaks

## View results

```bash
# List recent browser probes
oackctl browser-probes list --team <TEAM_ID> --monitor <MONITOR_ID>

# Get detailed results for a specific probe
oackctl browser-probes get --team <TEAM_ID> --monitor <MONITOR_ID> <PROBE_ID>

# Download screenshot
oackctl browser-probes screenshot --team <TEAM_ID> --monitor <MONITOR_ID> <PROBE_ID>

# Download HAR file (open in Chrome DevTools → Network → Import)
oackctl browser-probes har --team <TEAM_ID> --monitor <MONITOR_ID> <PROBE_ID>
```

## Script Security Model

Your script runs inside a 4-layer sandbox:
1. **Docker container** — isolated filesystem, network, and process namespace
2. **Network policy** — only outbound HTTPS allowed
3. **Code harness** — only `page` and `context` globals; no `require`, `process`, `fs`
4. **Resource limits** — CPU, memory, and time caps per execution
