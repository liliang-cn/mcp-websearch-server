package utils

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestRetryWithBackoff_Success(t *testing.T) {
	attempts := 0
	fn := func() error {
		attempts++
		if attempts < 2 {
			return errors.New("temporary error")
		}
		return nil
	}

	config := RetryConfig{
		MaxAttempts:  3,
		InitialDelay: 1 * time.Millisecond,
		MaxDelay:     10 * time.Millisecond,
		Multiplier:   2.0,
	}

	ctx := context.Background()
	err := RetryWithBackoff(ctx, config, fn)

	if err != nil {
		t.Errorf("expected success but got error: %v", err)
	}

	if attempts != 2 {
		t.Errorf("expected 2 attempts, got %d", attempts)
	}
}

func TestRetryWithBackoff_AllAttemptsFail(t *testing.T) {
	attempts := 0
	fn := func() error {
		attempts++
		return errors.New("persistent error")
	}

	config := RetryConfig{
		MaxAttempts:  3,
		InitialDelay: 1 * time.Millisecond,
		MaxDelay:     10 * time.Millisecond,
		Multiplier:   2.0,
	}

	ctx := context.Background()
	err := RetryWithBackoff(ctx, config, fn)

	if err == nil {
		t.Error("expected error but got success")
	}

	if attempts != 3 {
		t.Errorf("expected 3 attempts, got %d", attempts)
	}
}

func TestRetryWithBackoff_ContextCancellation(t *testing.T) {
	attempts := 0
	fn := func() error {
		attempts++
		return errors.New("error")
	}

	config := RetryConfig{
		MaxAttempts:  5,
		InitialDelay: 100 * time.Millisecond,
		MaxDelay:     1 * time.Second,
		Multiplier:   2.0,
	}

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	err := RetryWithBackoff(ctx, config, fn)

	if err != context.Canceled {
		t.Errorf("expected context.Canceled error, got %v", err)
	}

	if attempts > 2 {
		t.Errorf("expected at most 2 attempts before cancellation, got %d", attempts)
	}
}

func TestRetryWithBackoff_DelayBackoff(t *testing.T) {
	attempts := 0
	var delays []time.Duration
	lastTime := time.Now()

	fn := func() error {
		attempts++
		if attempts > 1 {
			delays = append(delays, time.Since(lastTime))
		}
		lastTime = time.Now()
		return errors.New("error")
	}

	config := RetryConfig{
		MaxAttempts:  3,
		InitialDelay: 10 * time.Millisecond,
		MaxDelay:     100 * time.Millisecond,
		Multiplier:   2.0,
	}

	ctx := context.Background()
	_ = RetryWithBackoff(ctx, config, fn)

	if len(delays) != 2 {
		t.Fatalf("expected 2 delays, got %d", len(delays))
	}

	if delays[0] < 10*time.Millisecond {
		t.Errorf("first delay too short: %v", delays[0])
	}

	if delays[1] < 20*time.Millisecond {
		t.Errorf("second delay too short: %v", delays[1])
	}
}

func TestDefaultRetryConfig(t *testing.T) {
	config := DefaultRetryConfig()

	if config.MaxAttempts != 3 {
		t.Errorf("expected MaxAttempts=3, got %d", config.MaxAttempts)
	}

	if config.InitialDelay != 1*time.Second {
		t.Errorf("expected InitialDelay=1s, got %v", config.InitialDelay)
	}

	if config.MaxDelay != 10*time.Second {
		t.Errorf("expected MaxDelay=10s, got %v", config.MaxDelay)
	}

	if config.Multiplier != 2.0 {
		t.Errorf("expected Multiplier=2.0, got %f", config.Multiplier)
	}
}
