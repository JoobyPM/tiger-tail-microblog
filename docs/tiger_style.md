# Tiger Style (Go)

> This is an adaptation of the
> [**Tiger Style**](https://github.com/tigerbeetle/tigerbeetle/blob/main/docs/TIGER_STYLE.md) guide for the Go (`golang`)
> ecosystem. While the spirit of the original document remains—focusing on safety, performance, and developer
> experience—the specifics have been updated to reflect Go’s type system, standard library, and concurrency model.

## The Essence Of Style

> “There are three things extremely hard: steel, a diamond, and to know one's self.” — Benjamin Franklin

Our coding style is evolving, shaped by the interplay of engineering and art, numbers and human intuition. We blend
reason, experience, first principles, and knowledge—striving for precision, elegance, and a little poetry in the code we
write. Just like a good beat or a rare groove, it’s about finding that flow. The best is yet to come.

## Why Have Style?

Another word for style is design.

> “Design is not just what it looks like and feels like. Design is how it works.” — Steve Jobs

Our design goals: **safety**, **performance**, and **developer experience**—in that order. Readability is important, but
it’s only table stakes. The deeper question is whether the code contributes to or detracts from safety, performance, or
developer experience. That’s why we have a style guide.

> “… in programming, style is not something to pursue directly. Style is necessary only where understanding is missing.”
> — [Let Over Lambda](https://letoverlambda.com/index.cl/guest/chap1.html)

Below, we’ll explore how we apply these design goals to coding style, from simplicity and elegance to handling technical
debt.

---

## On Simplicity And Elegance

Simplicity is not a compromise. In fact, it’s the key to unifying safety, performance, and developer experience—the
“super idea.” Achieving real simplicity takes multiple passes, thoughtful design, and sometimes the willingness to
“throw one away” to rebuild something far more coherent.

> “Simplicity and elegance are unpopular because they require hard work and discipline to achieve.” — Edsger Dijkstra

One day or week of extra thought in design can save you months in production. The payoff is huge.

## Technical Debt

**Zero technical debt** is our policy. Because once code is shipped, mistakes get exponentially more expensive to fix—
and good work builds momentum.

What we do ship might lack certain features, but it won’t compromise our design goals. This is how we make steady,
incremental progress.

---

## Safety

> “The rules act like the seat-belt in your car: initially they are perhaps a little uncomfortable, but after a while
> their use becomes second-nature and not using them becomes unimaginable.” — Gerard J. Holzmann

Go is memory-managed in the sense that the garbage collector handles allocation and freeing, but safety issues can still
arise through unchecked logic, concurrency hazards (race conditions), or unverified data. Below are guidelines adapted
from NASA’s [Power of Ten Rules](https://spinroot.com/gerard/pdf/P10.pdf) and other proven safety practices:

1. **Keep Control Flow Simple**
   - Avoid deeply nested or convoluted control flows. Clear loops and explicit conditions are safer than “clever”
     one-liners.
   - Use recursion sparingly—Go allows it, but deep recursion is less common. Large call stacks can lead to performance
     issues or unexpected logic problems.

2. **Set Limits Everywhere**
   - Even though the garbage collector manages memory, set clear upper bounds on loops, channel sizes, request payloads,
     etc. This helps prevent resource exhaustion.

3. **Use Strict Types**
   - Go is statically typed. Favor well-defined types over `interface{}` or generics that obscure constraints.
   - Use domain-specific types where clarity or safety matters (e.g., a `UserID` type instead of a bare `string`).

4. **Assert What You Expect**
   - Validate function inputs. Use error returns or panics (rarely, but intentionally) for genuinely impossible states.
   - For critical invariants, consider small guard functions or checks. Splitting assertions can make it clearer which
     part failed:
     ```go
     if a == nil {
       return errors.New("a is nil")
     }
     if b < 0 {
       return errors.New("b is negative")
     }
     ```
     is often clearer than
     ```go
     if a == nil || b < 0 {
       return errors.New("bad input")
     }
     ```

5. **Pair Checks**
   - If data is validated before writing to a database or external system, re-validate it after reading or retrieval.
   - This helps catch corruption or domain errors quickly.

6. **Static Analysis and Tooling**
   - Use `go vet`, `golangci-lint`, or other static analyzers to catch common pitfalls. They can flag unused variables,
     unreachable code, race conditions, etc.

7. **Memory and Concurrency**
   - Even with garbage collection, watch for references that can cause leaks (e.g., storing large slices in global
     variables).
   - Concurrency hazards: always guard shared data with mutexes or channels. Use `go run -race` in testing to detect
     race conditions.

8. **Keep Functions Short**
   - A rough limit: ~50–70 lines of _actual code_ per function. This reduces complexity and improves readability.
   - Decompose large logic blocks. Keep state manipulation in the main function; push simpler tasks into helper
     functions.

9. **Be Warned**
   - Resolve warnings and linter errors (where possible). A clean build with `go vet` and `golangci-lint` is often an
     indicator of good hygiene.

10. **Drive the Program; Don’t Let the Outside Drive You**
    - Wherever possible, structure your program to maintain control over how data is processed (e.g., using channels to
      batch requests) rather than letting external streams or large inputs push you into tight loops or unbounded
      memory usage.

Lastly, handle errors consistently. In distributed systems, many “minor” errors can lead to major catastrophes if they
aren’t handled or logged properly.

---

## Performance

> “The lack of back-of-the-envelope performance sketches is the root of all evil.” — Rivacindela Hudsoni

Even though Go compiles to native binaries and provides efficient concurrency primitives, it’s crucial to design for
performance from the start:

1. **Sketch Performance Early**
   - Before coding, estimate how much CPU, network, disk, and memory your design needs. Resist the temptation to “just
     optimize later”—the biggest wins come from design-level decisions.

2. **Slowest Resource First**
   - Network is usually the bottleneck, then disk, memory, and CPU. Batching network I/O or database calls typically
     yields bigger gains than micro-optimizing CPU loops.

3. **Avoid Surprising the Compiler & Runtime**
   - Go’s compiler and runtime can optimize predictable code paths. Keep hot loops free from frequent allocations or
     type-morphing maps or slices.

4. **Separate Control and Data Planes**
   - Clearly separate logic (control plane) from data transformations (data plane). Minimize back-and-forth across
     multiple layers when in a performance-critical path.

5. **Explicit is Better**
   - Avoid “magical” reflection-heavy patterns unless truly necessary. Reflection can hamper clarity and performance.
   - Simpler code is easier to optimize and to reason about.

---

## Developer Experience

> “There are only two hard things in Computer Science: cache invalidation, naming things, and off-by-one errors.” — Phil
> Karlton

### Naming Things

- **Nouns & Verbs**
  - Great naming is great design. Prefer descriptive, precise words that reveal intent.
  - If it’s an action, use a verb (`CreateUser`, `FetchRecords`). If it’s a thing, use a noun (`UserAccount`,
    `LedgerRecord`).

- **Case Conventions**
  - Go uses **CamelCase** for exported names and **camelCase** for unexported. No underscores or dashes in function or
    variable names.
  - For constants, stick to CamelCase if exported (e.g., `MaxRetries`), or `maxRetries` if it’s not exported.
  - Filenames are typically lowercase with underscores as needed (`user_service.go`).

- **Avoid Overloading Terms**
  - Don’t reuse domain terms that mean something else within the same codebase or in overlapping functionalities.

- **Use Unit Suffixes**
  - `latencyMs`, `countLimit`, `heightPx`—knowing the measurement helps prevent mistakes.

- **Commit Messages**
  - Write _why_ something changed, not just _what_ changed.

- **Comment Well**
  - Go uses `//` for line comments. Use these to briefly describe tricky logic or domain-specific concepts.
  - For exported symbols, use standard Go doc comments starting immediately before the declaration:
    ```go
    // CreateUser sets up a new user account in the system.
    func CreateUser(...) { ... }
    ```

### Cache Invalidation

- **Keep Variables Single-Purpose**
  - Avoid reusing the same variable for multiple roles; it reduces confusion and keeps state straightforward.

- **In-Place Mindset**
  - Even with garbage collection, an “in-place” mindset helps avoid conceptual sprawl. Build objects once or use
    well-defined structs instead of reassigning large slices or maps for multiple purposes.

- **Shrink Scope**
  - Declare variables in the tightest scope needed (`if` blocks, short helper functions). This reduces the risk of
    accidental misuse.

- **Simplify Functions & Return Types**
  - Keep the number of return values manageable (Go allows multiple returns, but more than three or four can be
    confusing).
  - Avoid too many optional or “mixed” types.

- **Watch for Slice/Reference Pitfalls**
  - In Go, slicing arrays or referencing sub-slices can lead to subtle bugs if the underlying array changes. If you need
    an immutable view, consider copying.

### Off-By-One Errors

- **Index vs. Length**
  - The usual suspect: watch for `i < len(slice)`, and confirm that your loops terminate correctly.

- **Division & Rounding**
  - Go integer division truncs toward zero. If you need floating math or rounding, do it explicitly (e.g., `math.Floor`,
    `math.Ceil`, or type conversions).

---

## Style By The Numbers

- **Use `go fmt`**
  - Let `go fmt` format your code. This is idiomatic Go and one of the biggest perks of the language: a standardized
    style.

- **Indentation**
  - Go uses **tabs** by convention, where `go fmt` handles alignment automatically.

- **Line Length**
  - There’s no official limit, but aim for readability. Many teams stick to about 100 columns. Let your editor
    accommodating `go fmt` handle wrapping.

- **Braces**
  - Go requires braces for `if`, `for`, `switch`, etc. Unlike some other languages, single-line statements still need
    braces:
    ```go
    if err != nil {
      return err
    }
    ```

---

## Dependencies

We prefer **minimal dependencies** beyond the Go standard library. Every new module can add supply chain risk,
versioning conflicts, or extra complexity. We treat dependencies carefully, especially in lower-level or critical
systems.

---

## Tooling

> “The right tool for the job is often the tool you are already using—adding new tools has a higher cost than many
> people appreciate.” — John Carmack

Go has a rich toolchain out of the box:

- **Built-in Test Runner**  
  Use `go test ./...` for unit and integration tests.

- **Built-in Formatting**  
  `go fmt` (or `goimports`) for consistent code formatting.

- **Vet & Lint**  
  `go vet` to catch suspicious constructs; `golangci-lint` or `staticcheck` for more in-depth static analysis.

- **Go Modules**  
  Manage dependencies with `go mod init`, `go mod tidy`, etc. Keep `go.mod` and `go.sum` under version control.

Whenever you can, keep your code in idiomatic Go. That means straightforward concurrency via goroutines and channels,
simple package layouts, and small, composable modules.

---

## The Last Stage

Finally, have fun and keep iterating. We call it “TigerBeetle” because we aim for code that is both incredibly fast and
impressively small—run swiftly and lightly, with minimal overhead, in a wide world full of possibilities.

> “You don’t really suppose, do you, that all your adventures and escapes were managed by mere luck, just for your sole
> benefit? You are a very fine person… but you are only quite a little fellow in a wide world after all!”
>
> “Thank goodness!” said Bilbo laughing, and handed him the tobacco-jar.

Remember, code that is simpler is safer, and code that’s safer is usually more elegant. And elegant code is just more
fun to work with.

**Happy coding!**