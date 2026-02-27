package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	PasteCreations = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "pastebin_paste_creations_total",
			Help: "Total number of pastes created",
		},
		[]string{"visibility"},
	)

	PasteViews = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "pastebin_paste_views_total",
			Help: "Total number of paste views",
		},
		[]string{"type"},
	)

	PasteUpdates = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "pastebin_paste_updates_total",
			Help: "Total number of paste updates",
		},
	)

	PasteDeletions = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "pastebin_paste_deletions_total",
			Help: "Total number of paste deletions",
		},
	)

	UserRegistrations = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "pastebin_user_registrations_total",
			Help: "Total number of user registrations",
		},
	)

	UserLogins = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "pastebin_user_logins_total",
			Help: "Total number of user logins",
		},
		[]string{"status"},
	)

	AnalyticsEvents = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "pastebin_analytics_events_total",
			Help: "Total number of analytics events tracked",
		},
		[]string{"event_type"},
	)

	ActivePastes = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "pastebin_active_pastes",
			Help: "Current number of active (non-expired) pastes",
		},
	)

	PasteSizeBytes = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "pastebin_paste_size_bytes",
			Help:    "Size of paste content in bytes",
			Buckets: []float64{100, 500, 1000, 5000, 10000, 50000, 100000, 500000, 1000000},
		},
		[]string{"visibility"},
	)
)

func RecordPasteCreation(visibility string) {
	PasteCreations.WithLabelValues(visibility).Inc()
}

func RecordPasteView(viewType string) {
	PasteViews.WithLabelValues(viewType).Inc()
}

func RecordPasteUpdate() {
	PasteUpdates.Inc()
}

func RecordPasteDeletion() {
	PasteDeletions.Inc()
}

func RecordUserRegistration() {
	UserRegistrations.Inc()
}

func RecordUserLogin(status string) {
	UserLogins.WithLabelValues(status).Inc()
}

func RecordAnalyticsEvent(eventType string) {
	AnalyticsEvents.WithLabelValues(eventType).Inc()
}

func SetActivePastes(count float64) {
	ActivePastes.Set(count)
}

func RecordPasteSize(sizeBytes float64, visibility string) {
	PasteSizeBytes.WithLabelValues(visibility).Observe(sizeBytes)
}
