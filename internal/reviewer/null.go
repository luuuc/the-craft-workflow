package reviewer

// NullReviewer is a no-op reviewer for graceful degradation.
type NullReviewer struct{}

func (r *NullReviewer) Name() string {
	return NameNone
}

func (r *NullReviewer) Available() bool {
	return true
}

func (r *NullReviewer) Review(req ReviewRequest) (ReviewResponse, error) {
	return ReviewResponse{
		Content: `No external reviewer available.

Consider these questions yourself:
- Is this the smallest useful version?
- What could go wrong?
- What are you assuming that might not be true?
- Who else should weigh in before you proceed?

To enable AI review: set CRAFT_AI_API_KEY
To enable Council review: install council-cli`,
		Reviewer: NameNone,
	}, nil
}
