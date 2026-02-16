"use client";

import { motion } from "framer-motion";
import { CodeBlock } from "../../CodeBlock";
import { fadeInUp } from "../../animations";

export const StatsCommand = () => (
  <motion.article 
    initial="hidden"
    animate="visible"
    variants={fadeInUp}
    className="prose prose-invert max-w-none"
  >
    <motion.h1 
      variants={fadeInUp}
      className="font-mono text-3xl font-bold text-foreground mb-6"
    >
      mako stats
    </motion.h1>
    <p className="text-lg text-muted-foreground mb-8 leading-relaxed">
      View detailed usage statistics and command patterns
    </p>
      <div className="space-y-6">
        <div>
          <h3 className="font-mono text-base font-semibold text-foreground mb-3">Basic Statistics</h3>
          <CodeBlock>{`mako stats

# Output:
# Mako Statistics
# ─────────────────────────────────────
# Total commands executed:    1,234
# Success rate:                94.5%
# Failed commands:             68
# Average commands per day:    42
# Most productive day:         Monday (avg 58 commands)
# First command:               2024-01-15
# Last command:                2024-02-12 (5 minutes ago)`}</CodeBlock>
        </div>

        <div>
          <h3 className="font-mono text-base font-semibold text-foreground mb-3">Detailed Breakdown</h3>
          <CodeBlock>{`# Show detailed statistics
mako stats --verbose

# Most used commands
mako stats --top-commands 10

# Command frequency by hour
mako stats --by-hour

# Command frequency by day of week
mako stats --by-day

# Command frequency by directory
mako stats --by-directory`}</CodeBlock>
        </div>

        <div>
          <h3 className="font-mono text-base font-semibold text-foreground mb-3">Time-Based Analysis</h3>
          <CodeBlock>{`# Statistics for specific period
mako stats --since "1 week ago"
mako stats --since "2024-01-01"

# Compare periods
mako stats --compare "last week" "this week"

# Monthly summary
mako stats --monthly

# Year in review
mako stats --year 2024`}</CodeBlock>
        </div>

        <div>
          <h3 className="font-mono text-base font-semibold text-foreground mb-3">Pattern Analysis</h3>
          <CodeBlock>{`# Find command patterns
mako stats --patterns

# Output:
# Common Patterns Detected:
# 1. Git workflow: commit → push (89% of git operations)
# 2. Docker workflow: build → up → logs
# 3. npm workflow: install → test → start
# 4. Most efficient time: 2pm-4pm (highest success rate)

# Productivity insights
mako stats --insights

# Command similarity analysis
mako stats --duplicates`}</CodeBlock>
        </div>

        <div>
          <h3 className="font-mono text-base font-semibold text-foreground mb-3">Export Options</h3>
          <CodeBlock>{`# Export stats as JSON
mako stats --format json --output stats.json

# Export as CSV
mako stats --format csv --output stats.csv

# Generate HTML report
mako stats --format html --output report.html`}</CodeBlock>
        </div>

        <div>
          <h3 className="font-mono text-base font-semibold text-foreground mb-3">Edge Cases</h3>
          <div className="space-y-3">
            <div className="border-l-2 border-blue-500/50 pl-4 py-2 hover:bg-muted-foreground/5 transition-colors rounded-r">
              <h4 className="font-mono text-sm font-semibold text-foreground mb-2">Insufficient Data</h4>
              <CodeBlock>{`mako stats`}</CodeBlock>
              <p className="text-muted-foreground text-sm leading-relaxed mt-2">
                Output: Not enough data for meaningful statistics. Run at least 10 commands to see statistics.
              </p>
            </div>

            <div className="border-l-2 border-red-500/50 pl-4 py-2 hover:bg-muted-foreground/5 transition-colors rounded-r">
              <h4 className="font-mono text-sm font-semibold text-foreground mb-2">Invalid Date Range</h4>
              <CodeBlock>{`mako stats --since "invalid date"`}</CodeBlock>
              <p className="text-muted-foreground text-sm leading-relaxed mt-2">
                Error: Invalid date format. Use: YYYY-MM-DD or relative dates like "1 week ago"
              </p>
            </div>
          </div>
        </div>
      </div>
    </motion.article>
  );
