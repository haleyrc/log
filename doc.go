// Package log provides a minimal logging interface primarily aimed at enabling
// structured logging.
//
// This package borrows heavily from https://github.com/apex/log, but vastly
// simplifies the API in order to railroad users into a consistent usage. The
// other primary difference is that this package is aware of tags stored the
// provided context coming from the https://github.com/haleyrc/tag package. This
// makes it significantly easier to provide a consistent set of tags for all
// related log lines. By combining tags with one-off fields passed directly to
// the logging methods, you get consistent, structured logs that can be ingested
// by any modern logging platform that can ingest one-JSON-object-per-line logs
// e.g. DataDog, Google Logging and BigQuery, Amazon Athena, etc.
//
// # Handlers
//
// A Logger is only directly responsible for marshaling the different data
// sources (level, tags, msg, fields, etc.) into a Message, which can be thought
// of as a type representing a single "log line". It is up to the Logger's
// MessageHandler to do the actual formatting.
//
// In order to keep the API simple and enforce the requirements that drove the
// development of this package, we only provide two implementations of the
// MessageHandler interface: JSON and Memory.
//
// The JSON implementation is intended for use in actual applications and does
// what it says on the tin: formats Messages as single-line JSON blobs.
//
// Conversely, the Memory implementation does no formatting, but stores the
// received Messages unchanged in an in-memory cache. The only suitable use for
// this MessageHandler is in testing scenarios where it's important to verify
// the content of logs directly.
//
// Note, if neither of the provided MessageHandler implementations works for
// you, you can create and provide your own handler.
//
// # Levels
//
// While this package makes use of some of the standard syslog terminology for
// naming the supported log levels, it's important to note that this is not
// intended to be a syslog package. Notably, a number of syslog levels have been
// intentionally omitted. For a detailed discussion about the motivation behind
// this decision, see the package README.
//
// What's left are the following levels: debug, info, and error. The names of
// these levels are largely self-explanatory, but the brief rundown of intended
// use is as follows:
//
//   - The "debug" level is used for logs that are of primary use during initial
//     development and troubleshooting. These are not intended to be sent to
//     live, production logs.
//   - The "error" level is used to report any erroneous conditions. These are
//     almost always treated differently by log ingestion platforms (assuming
//     your platform is correctly setup to parse this format).
//   - The "info" level is used for all the rest of your logging needs. Best
//     practices here are largely specific to your organization, but this
//     package recommends logging all non-error "domain events" at this level.
//     For a more detailed discussion of this topic, see the README.
package log
