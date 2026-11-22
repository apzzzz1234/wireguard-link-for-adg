# Integration readiness review (2026-11-22)

Scope: upstream v1.0.11, commit 9138c84d2c5f900455ad973953037344a95520f8.
This is a source review and automated host verification, not a security certification or Android endurance certification. No production profiles, peer registrations, servers, or installed apps were changed.

## Decision

Suitable for a one-device pilot; do not replace all existing VPN clients yet. Preserve the upstream compatibility package name. Do not reuse an internal app's active WireGuard key for this client: simultaneous clients with the same peer identity can cause endpoint roaming and disruption. Import only the intended device's profile, with its other client stopped.

## Source observations

- ProfileStore uses Android Keystore AES-256-GCM, AtomicFile, and no-backup storage. Android backup is disabled. Explicit user exports contain private keys and are NOT encrypted archives; treat exported files as secrets.
- MegazordService accepts sensitive IPC only from its own UID or a UID owning com.adguard.android. It does not pin AdGuard's signing certificate; install AdGuard from its official distribution.
- SOCKS binds to loopback, requires generated random username/password, limits clients, and applies authentication/connection deadlines. It cannot be configured to bind all network interfaces through startProxy.
- startProxy sends matching AllowedIPs through the WireGuard userspace network stack; nonmatching destinations go directly. A failed matching route does not fall back to direct transport.
- No hardcoded telemetry/upload endpoint found in application/native source. Build dependency downloads are separate from runtime traffic. This is static evidence, not a packet-capture result.
- Native statistics whitelist counters/handshake timestamps rather than returning raw WireGuard configuration. Java exception logging omits error messages/config values.
- Recovery retries are capped at 60 seconds between attempts. Binder replacement invalidates stale readiness. Endpoint names are resolved again on network changes.

## Important limits

- Connected means local engine running, not remote handshake/service health.
- IncludedApplications / ExcludedApplications are applied by buildTun (standalone mode). In AdGuard integration, filtering/app exclusions must be checked in AdGuard; do not assume those profile fields apply.
- Domain routing selects the first resolved address. No multi-address fallback. DNS configured in the tunnel must be reachable using the tunnel stack; inspect DNS as well as AllowedIPs during onboarding.
- SOCKS UDP with an unspecified source port learns its first sender; loopback UDP does not have per-packet authentication.
- Android reboot, vendor background restrictions, AdGuard upgrades, phone network transitions, battery and throughput remain unverified here.
- Fork signing differs from upstream: a future APK signed by us cannot update the upstream APK in place without a compatible signing key. Do not silently uninstall or erase profiles.

## Automated evidence

- Windows Go 1.27.1: upstream native tests plus two new routing/bind regressions passed; go vet passed on upstream baseline.
- Added TestAllowedDestinationNeverFallsBackToDirect: deliberately reachable host listener, matching tunnel route with no usable peer; SOCKS fails and host listener receives no connection.
- Added TestProxyRejectsPublicBindAndMissingRoutes: public bind addresses, empty and malformed routes are rejected.
- UI translation check passed: 153 messages, 129 uses.
- Maven dependency checksums verified by upstream fetch tool.
- govulncheck v1.10.0 on Windows host build: 0 reachable symbol findings; 3 package and 45 module findings, including unused paths and Windows-specific findings. This is NOT an Android build vulnerability clearance. Dependency upgrades and Android-target dependency review remain pre-release work.
- Upstream baseline Linux race tests / Android four-ABI build: see linked CI run; results are recorded after completion.

## Device acceptance checklist

1. Disable the other system VPN and use one dedicated device profile; enable AdG integration and confirm actual WG handshake plus authenticated internal service access.
2. Confirm intended private destination takes WG and ordinary Internet takes direct egress, including DNS and IPv6 behavior.
3. Keep photograph backup exclusions explicit; compare synthetic transfer throughput with existing baseline, without modifying photo history.
4. Switch Wi-Fi/mobile, interrupt network, restore it, and restart AdGuard: check automatic recovery and no silent direct access to protected destinations.
5. Lock screen for 30–60 minutes, then overnight; verify service access, battery use and notifications. Test reboot separately.
6. Check ICC/KodBox/Termix login, WebSocket and file transfer. AdGuard HTTPS inspection exceptions may be required; do not disable certificate verification in the apps.
7. Only after the pilot passes, produce a separately signed fork release with SHA-256 and update instructions. No production APK is published by this review.
