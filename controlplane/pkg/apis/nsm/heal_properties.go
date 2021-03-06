package nsm

import (
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
	"time"
)

const (
	NsmdHealEnabled        = "NSMD_HEAL_ENABLED"      // Does healing is enabled or not
	NsmdHealDSTWaitTimeout = "NSMD_HEAL_DST_TIMEOUTs" // Wait timeout for DST in seconds
)

type HealTimeouts struct {
	HealTimeout          time.Duration
	HealCloseTimeout     time.Duration
	HealRequestTimeout   time.Duration
	HealDataplaneTimeout time.Duration

	// Total DST heal timeout is 20 seconds.
	HealDSTNSEWaitTimeout time.Duration
	HealDSTNSEWaitTick    time.Duration

	HealEnabled bool
}

func NewHealProperties() *HealTimeouts {
	values := &HealTimeouts{
		HealTimeout:          time.Minute * 1,
		HealCloseTimeout:     time.Second * 5,
		HealRequestTimeout:   time.Minute * 1,
		HealDataplaneTimeout: time.Minute * 1,

		// Total DST heal timeout is 20 seconds.
		HealDSTNSEWaitTimeout: time.Second * 30,       // Maximum time to wait for NSMD/NSE to re-appear
		HealDSTNSEWaitTick:    500 * time.Millisecond, // Wait timeout to appear of NSE
		HealEnabled:           true,
	}

	// Parse few Environment variables.
	if "false" == os.Getenv(NsmdHealEnabled) {
		values.HealEnabled = false
	}
	dstWaitTimeout := os.Getenv(NsmdHealDSTWaitTimeout)
	if len(dstWaitTimeout) > 0 {
		logrus.Infof("Override HealDSTWaitTimeout: %s", dstWaitTimeout)
		seconds, err := strconv.ParseInt(dstWaitTimeout, 10, 32)
		if err == nil {
			values.HealDSTNSEWaitTimeout = time.Second * time.Duration(seconds)
		} else {
			logrus.Errorf("Failed to parse DST wait timeout value... %v", err)
		}
	}
	return values
}
