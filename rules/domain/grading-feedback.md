# Domain: Grading Feedback

Domain rules for the grading response contract. The whole product hinges on
*useful* feedback when an answer is wrong; a bare correct/incorrect boolean is
not a feature, it is a regression.

## Rules

### Rule: GRADE-001

- Owner layer: Domain
- Domain: grading-feedback
- Stability: core
- Status: active
- Scope: every grading endpoint response
- Statement: Grading responses MUST follow the envelope `{ correct, expected, explanation_zh, grammar_point, suggested_next[] }`. When `correct = false`, `explanation_zh` MUST contain a non-empty corrective explanation in Traditional Chinese, and `grammar_point` MUST point to an existing row in the `grammar_point` table.
- Rationale: The user's stated requirement: "錯誤的話告訴我要怎麼進行調整". The contract enforces that we always answer that question; nothing else is shippable.
- Verification: Schema validator on the response; integration test that posts a deliberately-wrong answer and asserts `len(explanation_zh) > 0` and `grammar_point` resolves to a real row.
- Supersedes: N/A
- Superseded by: N/A

### Rule: GRADE-002

- Owner layer: Domain
- Domain: grading-feedback
- Stability: core
- Status: active
- Scope: deterministic grader (no LLM)
- Statement: The deterministic grader MUST source `explanation_zh` from a curated `feedback_template` table keyed by `(question_id, error_class)`. It MUST NOT generate text on the fly. Unknown `error_class` falls back to a generic template that names the grammar point and links to its lesson page.
- Rationale: Deterministic-path responses must be reproducible and offline-capable. Generating text on the fly inside the deterministic path would defeat the whole reason that path exists (zero-cost, zero-LLM).
- Verification: Unit test ensuring the deterministic handler does not import any LLM client; lint ensuring no `feedback_template` row has an empty body.
- Supersedes: N/A
- Superseded by: N/A

### Rule: GRADE-003

- Owner layer: Domain
- Domain: grading-feedback
- Stability: behavior
- Status: active
- Scope: LLM grading paths (server provider + local connector)
- Statement: When an LLM grades a free-form answer, the prompt MUST instruct the model to (a) identify the most likely grammar-point error class, (b) cite the user's exact substring that triggered the error, and (c) give the corrective rewrite. Responses missing any of (a)/(b)/(c) MUST be rejected by the validator and retried up to once before falling back to a generic explanation.
- Rationale: Generic "your answer is wrong, try again" feedback is the failure mode we are explicitly trying to avoid. The three-part contract is what makes the feedback actionable.
- Verification: Validator unit tests covering each missing-part case; metric `grading.llm.retry_rate` exposed for monitoring.
- Supersedes: N/A
- Superseded by: N/A
