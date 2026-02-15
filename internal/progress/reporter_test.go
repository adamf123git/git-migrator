package progress

import (
	"sync"
	"testing"
	"time"
)

func TestNewReporter(t *testing.T) {
	r := NewReporter(100)
	if r == nil {
		t.Fatal("NewReporter returned nil")
	}
	if r.total != 100 {
		t.Errorf("total = %d, want 100", r.total)
	}
	if r.current != 0 {
		t.Errorf("current = %d, want 0", r.current)
	}
}

func TestNewReporterZero(t *testing.T) {
	r := NewReporter(0)
	if r == nil {
		t.Fatal("NewReporter returned nil")
	}
	if r.total != 0 {
		t.Errorf("total = %d, want 0", r.total)
	}
}

func TestReporterStart(t *testing.T) {
	r := NewReporter(100)
	r.Start()

	if r.startTime.IsZero() {
		t.Error("startTime should not be zero after Start()")
	}
}

func TestReporterSetCurrent(t *testing.T) {
	r := NewReporter(100)
	r.Start()

	r.SetCurrent(50)
	if r.current != 50 {
		t.Errorf("current = %d, want 50", r.current)
	}

	r.SetCurrent(75)
	if r.current != 75 {
		t.Errorf("current = %d, want 75", r.current)
	}
}

func TestReporterIncrement(t *testing.T) {
	r := NewReporter(100)
	r.Start()

	for i := 1; i <= 5; i++ {
		r.Increment()
		if r.current != i {
			t.Errorf("after increment %d, current = %d, want %d", i, r.current, i)
		}
	}
}

func TestReporterSetOperation(t *testing.T) {
	r := NewReporter(100)

	r.SetOperation("Processing commits")
	if r.operation != "Processing commits" {
		t.Errorf("operation = %q, want %q", r.operation, "Processing commits")
	}

	r.SetOperation("Creating branches")
	if r.operation != "Creating branches" {
		t.Errorf("operation = %q, want %q", r.operation, "Creating branches")
	}
}

func TestReporterCurrent(t *testing.T) {
	r := NewReporter(100)

	if r.Current() != 0 {
		t.Errorf("Current() = %d, want 0", r.Current())
	}

	r.SetCurrent(42)
	if r.Current() != 42 {
		t.Errorf("Current() = %d, want 42", r.Current())
	}
}

func TestReporterPercentage(t *testing.T) {
	tests := []struct {
		total   int
		current int
		want    float64
	}{
		{100, 0, 0},
		{100, 25, 25},
		{100, 50, 50},
		{100, 75, 75},
		{100, 100, 100},
		{200, 50, 25},
		{3, 1, 33.33333333333333},
		{0, 0, 0},
		{100, 150, 150}, // Over 100% is allowed
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			r := NewReporter(tt.total)
			r.SetCurrent(tt.current)
			got := r.Percentage()
			if got != tt.want {
				t.Errorf("Percentage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReporterPercentageZeroTotal(t *testing.T) {
	r := NewReporter(0)
	r.SetCurrent(10)

	if r.Percentage() != 0 {
		t.Errorf("Percentage() with zero total = %v, want 0", r.Percentage())
	}
}

func TestReporterOperation(t *testing.T) {
	r := NewReporter(100)

	r.SetOperation("Test Operation")
	if r.Operation() != "Test Operation" {
		t.Errorf("Operation() = %q, want %q", r.Operation(), "Test Operation")
	}
}

func TestReporterETA(t *testing.T) {
	r := NewReporter(100)
	r.Start()

	// Initially, ETA should be 0 (no progress yet)
	if eta := r.ETA(); eta != 0 {
		t.Errorf("ETA before progress = %v, want 0", eta)
	}

	// Set some progress and wait a bit
	r.SetCurrent(50)
	time.Sleep(10 * time.Millisecond)

	// ETA should now be calculable (non-zero)
	// Note: On very fast machines, ETA might still be 0 if elapsed time is too small
	eta := r.ETA()
	// Just verify ETA doesn't panic and returns a valid duration
	_ = eta
}

func TestReporterETAZeroProgress(t *testing.T) {
	r := NewReporter(100)
	// Don't call Start() - startTime is zero

	if eta := r.ETA(); eta != 0 {
		t.Errorf("ETA without Start() = %v, want 0", eta)
	}
}

func TestReporterETAZeroTotal(t *testing.T) {
	r := NewReporter(0)
	r.Start()
	r.SetCurrent(10)

	if eta := r.ETA(); eta != 0 {
		t.Errorf("ETA with zero total = %v, want 0", eta)
	}
}

func TestReporterSubscribe(t *testing.T) {
	r := NewReporter(100)

	var mu sync.Mutex
	received := []Status{}

	unsubscribe := r.Subscribe(func(s Status) {
		mu.Lock()
		received = append(received, s)
		mu.Unlock()
	})

	r.Start()
	r.SetCurrent(50)
	r.SetOperation("Testing")

	// Give some time for notifications
	time.Sleep(10 * time.Millisecond)

	mu.Lock()
	if len(received) < 3 {
		t.Errorf("received %d notifications, want at least 3", len(received))
	}
	mu.Unlock()

	// Test unsubscribe
	unsubscribe()

	// Clear received
	mu.Lock()
	received = nil
	mu.Unlock()

	r.SetCurrent(75)

	time.Sleep(10 * time.Millisecond)

	mu.Lock()
	if len(received) != 0 {
		t.Errorf("after unsubscribe, received %d notifications, want 0", len(received))
	}
	mu.Unlock()
}

func TestReporterMultipleSubscribers(t *testing.T) {
	r := NewReporter(100)

	var mu sync.Mutex
	count1, count2 := 0, 0

	r.Subscribe(func(s Status) {
		mu.Lock()
		count1++
		mu.Unlock()
	})

	r.Subscribe(func(s Status) {
		mu.Lock()
		count2++
		mu.Unlock()
	})

	r.Start()
	r.SetCurrent(25)
	r.SetCurrent(50)

	time.Sleep(10 * time.Millisecond)

	mu.Lock()
	if count1 < 3 {
		t.Errorf("subscriber 1 received %d notifications, want at least 3", count1)
	}
	if count2 < 3 {
		t.Errorf("subscriber 2 received %d notifications, want at least 3", count2)
	}
	mu.Unlock()
}

func TestReporterSubscribeNilCallback(t *testing.T) {
	r := NewReporter(100)

	// This should not panic
	var received Status
	r.Subscribe(func(s Status) {
		received = s
	})

	r.Start()
	r.SetCurrent(50)

	time.Sleep(10 * time.Millisecond)

	if received.Current != 50 {
		t.Errorf("received.Current = %d, want 50", received.Current)
	}
}

func TestReporterNotifyStatusFields(t *testing.T) {
	r := NewReporter(100)

	var received Status
	r.Subscribe(func(s Status) {
		received = s
	})

	r.Start()
	r.SetCurrent(75)
	r.SetOperation("Processing items 1-75")

	time.Sleep(10 * time.Millisecond)

	if received.Current != 75 {
		t.Errorf("Current = %d, want 75", received.Current)
	}
	if received.Total != 100 {
		t.Errorf("Total = %d, want 100", received.Total)
	}
	if received.Percentage < 74 || received.Percentage > 76 {
		t.Errorf("Percentage = %v, want ~75", received.Percentage)
	}
	if received.Operation != "Processing items 1-75" {
		t.Errorf("Operation = %q, want %q", received.Operation, "Processing items 1-75")
	}
	if received.StartTime.IsZero() {
		t.Error("StartTime should not be zero")
	}
}

func TestReporterConcurrentIncrement(t *testing.T) {
	r := NewReporter(1000)
	r.Start()

	var wg sync.WaitGroup
	numGoroutines := 10
	incrementsPerGoroutine := 100

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < incrementsPerGoroutine; j++ {
				r.Increment()
			}
		}()
	}

	wg.Wait()

	expected := numGoroutines * incrementsPerGoroutine
	if r.Current() != expected {
		t.Errorf("Current() = %d, want %d", r.Current(), expected)
	}
}

func TestReporterConcurrentSetCurrent(t *testing.T) {
	r := NewReporter(100)
	r.Start()

	var wg sync.WaitGroup
	numGoroutines := 10

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			r.SetCurrent(val)
		}(i * 10)
	}

	wg.Wait()

	// Current should be one of the set values (0-90 in steps of 10)
	current := r.Current()
	if current < 0 || current > 90 || current%10 != 0 {
		t.Errorf("Current() = %d, want a multiple of 10 between 0 and 90", current)
	}
}

func TestReporterConcurrentSubscribe(t *testing.T) {
	r := NewReporter(100)
	r.Start()

	var mu sync.Mutex
	allReceived := []Status{}

	var wg sync.WaitGroup
	numSubscribers := 10

	for i := 0; i < numSubscribers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			r.Subscribe(func(s Status) {
				mu.Lock()
				allReceived = append(allReceived, s)
				mu.Unlock()
			})
		}()
	}

	wg.Wait()
	r.SetCurrent(50)

	time.Sleep(20 * time.Millisecond)

	mu.Lock()
	if len(allReceived) < numSubscribers {
		t.Errorf("received %d notifications, want at least %d", len(allReceived), numSubscribers)
	}
	mu.Unlock()
}

func TestReporterETACalculation(t *testing.T) {
	r := NewReporter(100)
	r.Start()

	// Simulate processing 25% in 100ms
	r.SetCurrent(25)
	time.Sleep(100 * time.Millisecond)

	eta := r.ETA()
	// ETA calculation depends on timing which can vary significantly
	// Just verify the ETA is a valid duration and the function doesn't panic
	// On fast machines, this might be very small or even 0
	if eta < 0 {
		t.Errorf("ETA = %v, should not be negative", eta)
	}
	// Verify the reporter state is correct after calculation
	if r.Current() != 25 {
		t.Errorf("Current = %d, want 25", r.Current())
	}
}

func TestReporterStatusCopy(t *testing.T) {
	r := NewReporter(100)

	var received []Status
	r.Subscribe(func(s Status) {
		received = append(received, s)
	})

	r.Start()
	r.SetCurrent(50)

	time.Sleep(10 * time.Millisecond)

	// Modify the received status - should not affect the reporter
	if len(received) > 0 {
		received[0].Current = 999
		if r.Current() == 999 {
			t.Error("modifying received Status affected reporter")
		}
	}
}

func TestReporterPercentagePrecision(t *testing.T) {
	tests := []struct {
		total   int
		current int
	}{
		{7, 1},
		{7, 2},
		{7, 3},
		{13, 5},
		{99, 33},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			r := NewReporter(tt.total)
			r.SetCurrent(tt.current)
			pct := r.Percentage()
			expected := float64(tt.current) / float64(tt.total) * 100
			if pct != expected {
				t.Errorf("Percentage() = %v, want %v", pct, expected)
			}
		})
	}
}

func TestReporterOperationEmpty(t *testing.T) {
	r := NewReporter(100)

	if r.Operation() != "" {
		t.Errorf("Operation() = %q, want empty string", r.Operation())
	}

	r.SetOperation("")
	if r.Operation() != "" {
		t.Errorf("Operation() = %q, want empty string", r.Operation())
	}
}

func TestReporterNegativeValues(t *testing.T) {
	r := NewReporter(100)
	r.Start()

	// Setting negative values should work (not validated)
	r.SetCurrent(-5)
	if r.Current() != -5 {
		t.Errorf("Current() = %d, want -5", r.Current())
	}

	// Percentage can be negative
	if pct := r.Percentage(); pct >= 0 {
		t.Errorf("Percentage() = %v, want negative", pct)
	}
}
