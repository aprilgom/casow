# casow 2차 포팅 ROI 산정

## 기준

- 원본: `/Users/aprilgom/casow/kasuari`
- 대상: `/Users/aprilgom/casow/casow`
- 원칙: 테스트를 먼저 포팅하고, 그 후 코드 포팅
- 현재 상태: foundation PR #1이 `main`에 merge됨. 외부 integration 테스트 중 `removal`, `quadrilateral`, `traits`는 Go 테스트로 반영됨.

## 이미 닫힌 범위

- public value API: `Variable`, `Term`, `Expression`, `Strength`, `Constraint`, relation/error 값 타입
- internal tableau: `symbol`, `row`
- solver API: `AddConstraint`, `AddConstraints`, `RemoveConstraint`, `AddEditVariable`, `RemoveEditVariable`, `SuggestValue`, `FetchChanges`, `GetValue`, `Reset`
- failure rollback: failed `AddConstraint` rollback 및 edit/change tracking 회귀 테스트
- zero change tracking: 마지막 제약 또는 edit variable 제거 후 `FetchChanges`가 `x=0`을 보고하는 동작
- upstream integration parity: `tests/removal.rs`, `tests/quadrilateral.rs`, `tests/traits.rs`

## ROI 순 다음 작업

### A. README/lib usage scenario parity test

ROI: 높음

원본 `kasuari/src/lib.rs`의 horizontal boxes 예제는 solver의 실제 사용 흐름을 가장 잘 압축한다. 현재 `usage_test.go`는 300/75 width 변화를 확인하지만, 원본 문서의 마지막 ratio constraint 추가 단계는 아직 없다.

테스트 먼저 포팅:

- window width 300 제안 후 known values 확인
- width 75 제안 후 required constraints 보존 확인
- ratio constraint 추가 후 `box1.right == 25`, `box2.left == 25`에 해당하는 변화 확인

예상 구현:

- 구현 변경 없음이 기대됨
- 실패 시 weak/medium strength ordering 또는 expression arithmetic 쪽 결함 가능성이 높음

병렬성:

- 독립적. solver behavior 테스트만 추가하므로 다른 API 테스트와 병렬 진행 가능

### B. change tracking lifecycle parity 확장

ROI: 높음

foundation에서 zeroed tracking을 보강했지만, upstream 문서가 강조하는 `FetchChanges` semantics는 더 체계적으로 고정할 가치가 있다.

테스트 먼저 포팅:

- initial zero value는 첫 `FetchChanges`에서 보고되지 않음
- nonzero initial solved value는 첫 `FetchChanges`에서 보고됨
- 같은 값으로 다시 suggest하면 변경이 보고되지 않음
- 여러 edit variables의 변경이 순서와 무관하게 모두 보고됨

예상 구현:

- 대부분 구현 변경 없음이 기대됨
- 실패 시 `changed` / `shouldClearChanges` lifecycle 수정 필요

병렬성:

- A와 독립적. `solver_test.go`를 같이 만질 수 있어 충돌 가능성은 있으나 테스트 함수 단위로 작게 나누면 병렬 가능

### C. public API parity gap: Go helper surface

ROI: 중간

Rust는 operator overloading으로 `Variable + Term`, `Term + Expression`, `PartialConstraint | rhs` 등을 제공한다. Go는 같은 문법을 포팅할 수 없으므로, Go helper API가 같은 표현력을 갖는지 테스트로 고정해야 한다.

테스트 먼저 포팅:

- `Var(x).PlusExpression(...)`, `TermFromVariable(x)`, `Expression` arithmetic 조합으로 원본 operator tests의 핵심 표현을 재현
- `NewConstraint(lhs, op, rhs, strength)`가 rhs `Variable`, `Term`, `Expression`, `constant` 역할을 모두 대체하는지 확인
- relation helper `EQ/LE/GE`가 weighted relation을 명확히 구성하는지 확인

예상 구현:

- 필요 시 작은 helper 추가 가능
- API 추가는 README 예제와 맞춰야 함

병렬성:

- solver internals와 독립적. A/B와 병렬 가능

### D. upstream internal arithmetic tests completeness

ROI: 중간

`variable.rs`, `term.rs`, `expression.rs`, `strength.rs`의 Rust 단위 테스트 중 Go에서 문법상 직접 포팅되지 않은 연산자 테스트가 많다. 현재 Go 테스트는 핵심 동작을 넓게 덮지만, upstream parity matrix로 빠진 케이스를 명시하면 이후 포팅 누락을 줄일 수 있다.

테스트 먼저 포팅:

- term/expression add/sub/mul/div 조합의 table test 보강
- strength create/add/sub/mul/div clamp table을 upstream case와 1:1로 대조
- relation/constraint construction table 보강

예상 구현:

- 구현 변경 없음이 기대됨
- 실패 시 value type helper의 arithmetic semantics 수정

병렬성:

- 파일별로 독립적이라 병렬 진행에 적합

### E. removal/pivot behavior stress cases

ROI: 중간

현재 removal 테스트는 단순 replacement와 zero change를 포함한다. 다음 solver 리스크는 remove가 external row pivot을 유발하는 복합 케이스다.

테스트 먼저 포팅:

- 여러 required/weak constraint를 추가한 뒤 중간 constraint 제거
- 제거 후 still-satisfiable required constraints와 expected values 확인
- `FetchChanges`가 stale value를 내지 않는지 확인

예상 구현:

- 실패 시 `getMarkerLeavingRow`, `removeConstraintEffects`, objective update 쪽 수정 가능

병렬성:

- A/B와 같은 `solver_test.go`를 만질 가능성이 커서 단독 진행 권장

## 추천 병렬 배치

첫 배치:

- A. README/lib usage scenario parity test
- C. public API parity gap: Go helper surface
- D. upstream internal arithmetic tests completeness

이 셋은 서로 쓰기 범위가 비교적 분리되어 있다. A는 `usage_test.go`, C는 `expression_test.go`/`constraint_test.go`/`relation_test.go`, D는 `term_test.go`/`strength_test.go` 중심으로 나누면 충돌이 작다.

두 번째 배치:

- B. change tracking lifecycle parity 확장
- E. removal/pivot behavior stress cases

둘 다 `solver_test.go` 중심이라 같은 시점에 병렬 진행하면 충돌 가능성이 있다. B를 먼저 닫고 E로 넘어가는 편이 안전하다.

## 다음 추천 작업

`subagent(low)` 병렬 작업으로 첫 배치 A + C + D를 진행한다. 각 subagent는 테스트를 먼저 추가하고, 구현 변경은 실패한 테스트를 통과시키는 최소 범위로 제한한다.
