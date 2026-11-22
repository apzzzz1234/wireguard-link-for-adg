# Fork maintenance

This fork retains upstream attribution and the compatibility package ID. Runtime code is unchanged in the initial review; the fork adds isolation regression tests and an integration readiness report.

- Upstream: https://github.com/david419kr/wireguard-link-for-adg
- Review: [2026-11-22](ICC-REVIEW-20261122.md)
- Track upstream changes with an upstream remote; review diffs before merging. Do not auto-install releases on production devices.
- Keep compatibility protocol changes, dependency upgrades and UI changes separate when possible.
- Run the full Build and test workflow for code/dependency changes. It checks Linux race tests and four-ABI Android packaging plus Java host checks.
- Keep all personal configurations, keys, signing material and server credentials out of this public fork. Tests use synthetic local addresses/keys only.
- Before distributing a fork APK, establish a private stable signing key and document migration from the upstream signer. This review has not published an installable fork APK.
- Before adopting an AdGuard or Android update, repeat the device acceptance checklist. The integration protocol is unofficial and may change independently.
