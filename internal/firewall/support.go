package firewall

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"os/exec"
	"strings"

	"github.com/qdm12/golibs/command"
)

var (
	ErrNetAdminMissing      = errors.New("NET_ADMIN capability is missing")
	ErrTestRuleCleanup      = errors.New("failed cleaning up test rule")
	ErrInputPolicyNotFound  = errors.New("input policy not found")
	ErrIPTablesNotSupported = errors.New("no iptables supported found")
)

func checkIptablesSupport(ctx context.Context, runner command.Runner,
	iptablesPathsToTry ...string) (iptablesPath string, err error) {
	var lastUnsupportedMessage string
	for _, pathToTest := range iptablesPathsToTry {
		ok, unsupportedMessage, err := testIptablesPath(ctx, pathToTest, runner)
		if err != nil {
			return "", fmt.Errorf("for %s: %w", pathToTest, err)
		} else if ok {
			iptablesPath = pathToTest
			break
		}

		lastUnsupportedMessage = unsupportedMessage
	}

	if iptablesPath == "" { // all iptables to try failed
		return "", fmt.Errorf("%w: from %s: last error is: %s",
			ErrIPTablesNotSupported, strings.Join(iptablesPathsToTry, ", "),
			lastUnsupportedMessage)
	}

	return iptablesPath, nil
}

func testIptablesPath(ctx context.Context, path string,
	runner command.Runner) (ok bool, unsupportedMessage string,
	criticalErr error) {
	// Just listing iptables rules often work but we need
	// to modify them to ensure we can support the iptables
	// being tested.

	// Append a test rule with a random interface name to the OUTPUT table.
	// This should not affect existing rules or the network traffic.
	testInterfaceName := randomInterfaceName()
	cmd := exec.CommandContext(ctx, path,
		"-A", "OUTPUT", "-o", testInterfaceName, "-j", "DROP")
	output, err := runner.Run(cmd)
	if err != nil {
		if isPermissionDenied(output) {
			// If the error is related to a denied permission,
			// return an error describing what to do from an end-user
			// perspective. This is a critical error and likely
			// applies to all iptables.
			criticalErr = fmt.Errorf("%w: %s", ErrNetAdminMissing, output)
			return false, "", criticalErr
		}
		unsupportedMessage = fmt.Sprintf("%s (%s)", output, err)
		return false, unsupportedMessage, nil
	}

	// Remove the random rule added previously for test.
	cmd = exec.CommandContext(ctx, path,
		"-D", "OUTPUT", "-o", testInterfaceName, "-j", "DROP")
	output, err = runner.Run(cmd)
	if err != nil {
		// this is a critical error, we want to make sure our test rule gets removed.
		criticalErr = fmt.Errorf("%w: %s (%s)", ErrTestRuleCleanup, output, err)
		return false, "", criticalErr
	}

	// Set policy as the existing policy so no mutation is done.
	// This is an extra check for some buggy kernels where setting the policy
	// does not work.
	cmd = exec.CommandContext(ctx, path, "-L", "INPUT")
	output, err = runner.Run(cmd)
	if err != nil {
		if isPermissionDenied(output) {
			criticalErr = fmt.Errorf("%w: %s", ErrNetAdminMissing, output)
			return false, "", criticalErr
		}
		unsupportedMessage = fmt.Sprintf("%s (%s)", output, err)
		return false, unsupportedMessage, nil
	}

	var inputPolicy string
	for _, line := range strings.Split(output, "\n") {
		inputPolicy, ok = extractInputPolicy(line)
		if ok {
			break
		}
	}

	if inputPolicy == "" {
		criticalErr = fmt.Errorf("%w: in INPUT rules: %s", ErrInputPolicyNotFound, output)
		return false, "", criticalErr
	}

	// Set the policy for the INPUT table to the existing policy found.
	cmd = exec.CommandContext(ctx, path, "--policy", "INPUT", inputPolicy)
	output, err = runner.Run(cmd)
	if err != nil {
		if isPermissionDenied(output) {
			criticalErr = fmt.Errorf("%w: %s", ErrNetAdminMissing, output)
			return false, "", criticalErr
		}
		unsupportedMessage = fmt.Sprintf("%s (%s)", output, err)
		return false, unsupportedMessage, nil
	}

	return true, "", nil // success
}

func isPermissionDenied(errMessage string) (ok bool) {
	const permissionDeniedString = "Permission denied (you must be root)"
	return strings.Contains(errMessage, permissionDeniedString)
}

func extractInputPolicy(line string) (policy string, ok bool) {
	const prefixToFind = "Chain INPUT (policy "
	i := strings.Index(line, prefixToFind)
	if i == -1 {
		return "", false
	}

	startIndex := i + len(prefixToFind)
	endIndex := strings.Index(line, ")")
	if endIndex < 0 {
		return "", false
	}

	policy = line[startIndex:endIndex]
	policy = strings.TrimSpace(policy)
	if policy == "" {
		return "", false
	}

	return policy, true
}

func randomInterfaceName() (interfaceName string) {
	const size = 15
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyz0123456789")
	b := make([]rune, size)
	for i := range b {
		letterIndex := rand.Intn(len(letterRunes)) //nolint:gosec
		b[i] = letterRunes[letterIndex]
	}
	return string(b)
}
