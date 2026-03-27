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

Create `test-local.js` — a thin wrapper that sets up the same environment your
script will have in the Oack sandbox:

```js
// test-local.js — run check.js locally with Playwright
const { chromium } = require("playwright");

(async () => {
  const browser = await chromium.launch({ headless: true });
  const browserContext = await browser.newContext({
    viewport: { width: 1920, height: 1080 },
  });
  const page = await browserContext.newPage();

  // Build the context object with env vars (same as Oack's script_env).
  const context = {
    LOGIN_EMAIL: process.env.LOGIN_EMAIL || "test@example.com",
    LOGIN_PASSWORD: process.env.LOGIN_PASSWORD || "s3cret",
  };

  try {
    // Load and run the check script — same call signature as the sandbox.
    const checkFn = require("./check.js");
    await checkFn(page, context);

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
