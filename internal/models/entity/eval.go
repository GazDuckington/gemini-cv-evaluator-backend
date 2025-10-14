package entity

type Rubric struct {
	Name        string  `json:"name"`
	Weight      float64 `json:"weight"`
	Description string  `json:"description"`
	Scale       string  `json:"scale"`
}

type EvaluationRubrics struct {
	CV      []Rubric `json:"cv"`
	Project []Rubric `json:"project"`
}

// NewDefaultRubrics returns the standard “Rubrics Cube” set.
func NewDefaultRubrics() EvaluationRubrics {
	return EvaluationRubrics{
		CV: []Rubric{
			{
				Name:        "Technical Skills Match",
				Weight:      0.40,
				Description: "Alignment with backend, databases, APIs, cloud, AI/LLM.",
				Scale:       "1=Irrelevant, 2=Few overlaps, 3=Partial, 4=Strong, 5=Excellent+AI/LLM",
			},
			{
				Name:        "Experience Level",
				Weight:      0.25,
				Description: "Years of experience and project complexity.",
				Scale:       "1=<1yr, 2=1–2yrs, 3=2–3yrs, 4=3–4yrs, 5=5+yrs",
			},
			{
				Name:        "Relevant Achievements",
				Weight:      0.20,
				Description: "Impact of past work.",
				Scale:       "1=None, 5=Major measurable impact",
			},
			{
				Name:        "Cultural/Collaboration Fit",
				Weight:      0.15,
				Description: "Communication, teamwork, learning mindset.",
				Scale:       "1=None, 5=Excellent and well-demonstrated",
			},
		},
		Project: []Rubric{
			{
				Name:        "Correctness",
				Weight:      0.30,
				Description: "Prompt chaining, RAG context injection.",
				Scale:       "1=None, 5=Fully correct + thoughtful",
			},
			{
				Name:        "Code Quality & Structure",
				Weight:      0.25,
				Description: "Clean, modular, tested.",
				Scale:       "1=Poor, 5=Excellent+strong tests",
			},
			{
				Name:        "Resilience & Error Handling",
				Weight:      0.20,
				Description: "API failures, retries, robustness.",
				Scale:       "1=None, 5=Production-ready",
			},
			{
				Name:        "Documentation & Explanation",
				Weight:      0.15,
				Description: "README clarity, setup, trade-offs.",
				Scale:       "1=Missing, 5=Excellent",
			},
			{
				Name:        "Creativity / Bonus",
				Weight:      0.10,
				Description: "Extra features beyond requirements.",
				Scale:       "1=None, 5=Outstanding creativity",
			},
		},
	}
}
