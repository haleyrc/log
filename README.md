# log

[![Build Status](https://github.com/haleyrc/log/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/haleyrc/log/actions?query=branch%3Amain)
[![GoDoc](https://pkg.go.dev/badge/github.com/haleyrc/log?status.svg)](https://pkg.go.dev/github.com/haleyrc/log?tab=doc)

A minimal API for structured logging.

## Install

```
$ go get -u github.com/haleyrc/log
```

## Usage

There are two ways to use the `log` package: using the default logger and using a custom logger. Which method to use is highly dependent on your own use-case and logging sensibilities.

For simple cases or for applications where loggers aren't passed around as explicit dependencies, it may be enough to use the top-level methods. These all use a "default logger", which is initialized using the `JSON` handler and with debug logging disabled. You can customize this default using the `SetDebug` and `SetHandler` functions, but be aware that the values you set will change for _all_ consumers. Assuming the defaults work for you, you can get going immediately:

```go
log.Info(ctx, "this is a test", log.F{
    "user_id": 123,
    "role":    "admin",
})
log.Error(ctx, "this is a test", log.F{
    "user_id": 123,
    "role":    "admin",
})
```

This results in log lines that look like:

```json
{"level":"INFO","timestamp":"2023-01-31T21:22:15.778736-05:00","message":"this is a test","tags":{},"fields":{"role":"admin","user_id":123}}
{"level":"ERROR","timestamp":"2023-01-31T21:22:15.77877-05:00","message":"this is a test","tags":{},"fields":{"role":"admin","user_id":123}}
```

If you want more control over your logger or you want to pass it around a la the recommendations from [Peter Bourgon](https://peter.bourgon.org/go-best-practices-2016/#logging-and-instrumentation), you can do that too.

The equivalent version of the example above would look like this:

```go
logger := log.NewJSONLogger(os.Stdout, nil)
logger.Info(ctx, "this is a test", log.F{
    "user_id": 123,
    "role":    "admin",
})
logger.Error(ctx, "this is a test", log.F{
    "user_id": 123,
    "role":    "admin",
})
```

That's about all there is to the library. If you're interested in why some things you're familiar with were included by others weren't, take a look at the section on [Motivation](#motivation). You might also be interested in my own personal suggestions for logging [Best Practices](#best-practices).

## Motivation

This package was designed entirely around things that have worked well for me writing, maintaining, and toubleshooting backend APIs. As with anything in software, there are thousand different ways to approach this problem and if you find that this package doesn't fit your use-case, I highly encourage you to check out any of the other high-quality logging packages out there. Here are some of my favorites:

<dl>
    <dt>https://github.com/apex/log</dt>
    <dd>This library is very obviously the inspiration for most of what you'll see in this package, but with way more bells and whistles.</dd>
    <dt>https://github.com/go-kit/kit</dt>
    <dd>The venerable go-kit suite of tools provides a great logger as part of its bag of tricks. If you're already using go-kit for other functionality, this is probably the first option you should check out.</dd>
    <dt>https://pkg.go.dev/log</dt>
    <dd>Sometimes, the standard library is all you need. Start here if you're just getting started and aren't sure what your requirements are.</dd>
</dl>

If you're still here, the following sections will help explain what drove the design of this package.

### Levels

If you look at [Wikipedia](https://en.wikipedia.org/wiki/Syslog#Severity_level), you'll see that the syslog standard includes a total of eight severity levels. By contrast, this package only provides three (and only two if you're running with debug off such as in a production environment). It may seem strange to limit the space of available options to this extent, but my guess is that anyone who has spent a lot of time trying to categorize log messages according to the full syslog standard has run into the same issues that I have. Specifically, **it's not at all clear which severity level to use for any given message**.

Sure, you can develop some standards and maybe even get pretty close, but I would wager the time spent to get to that point will almost certainly be lost. As with many things, there are diminishing returns involved with classifying log messages.

So instead of providing a Swiss Army knife with 100 tools nobody knows how to use, this package takes the opposite approach. In a production setting, we only expose two different levels: info and error. Compared to syslog, it's remarkably easy to determine which level to use. If you're logging an aberrant condition, use error. Otherwise, use info. Outside of that, a debug level is provided for the kinds of information that you often want to see while developing or troubleshooting, but that shouldn't print in a live scenario.

> **N.B.** There's nothing stopping you from immediately enabling debug output in a production scenario, but I would strongly advise against it. Apart from the fact that you'll quickly find yourself holding a huge infrastructure bill, the amount of noise generated will make it unnecessarily difficult to parse your logs.

After a number of years diagnosing production errors starting from the logs, I've found this to be simple enough to use correctly while still making a meaningful distinction between normal logs and aberrant conditions.

## Best Practices

In the following sections, I'll lay out what I consider to be the "best practices" for logging. Note that these recommendations come largely from writing backend API servers, so if you're writing CLI tools or some other type of application these suggestions may not apply. That said, I think it should be pretty simple to modify them to meet your needs.

### What to log

This section is mostly going to focus on what to log at the _info_ level, because it's usually pretty obvious when to use the error level (the rule for the debug level is just "log whatever you need to develop/troubleshoot so discussion on what to log is omitted). Note that most applications will have a smattering of logs that are more like instrumentation (e.g. Apache-style logs, timing logs, etc.). I also won't talk about these here since the inclusion/exclusion of these types of logs is so heavily dependendent on your architecture and supporting tooling.

With that in mind, what _should_ you log at the info level? Well, the glib answer is log whatever you need to. What I mean by that is this: if any recommendation for what to log doesn't include enough for you to do _your_ job in _your_ environment, log more. It's certainly possible to log too much, but that's almost never the problem people are facing. For instance, if you're logging a lot and you believe the logs should be useful but you're finding it too noisy to look through...you shouldn't assume you're logging too much. Instead, this is usually a result of either not using structured logging, not considering how you'll be searching your logs when you write them, or not setting up your logging platform correctly. I would address these issues before I started removing logging from my application.

> **N.B.** If your primary concern with logging too much is cost, you have a different issue entirely. Logs are your most useful tool for maintaining applications and you should cut costs just about anywhere else first. You might also be on the wrong platform.

I still haven't really made any suggestions in terms of what to log, so let's do that. I can't remember exactly where I read it (though I suspect it was from Peter Bourgon), but I've always liked the suggestion to "log domain events". This is both nebulous and immediately understandable (at least it is if your experience is anything like mine). Put another way, you should probably log anything that has a side effect in your system. Adding a new user to the database? Sending an email? Log all of them. Critically, log all of them with enough information to troubleshoot any issues that may arise **whether the event was successful or not**. More times than I can count, troubleshooting efforts were stymied because an operation succeeded...incorrectly. Don't make the mistake of believing that if a request goes off without a hitch, you won't have to troubleshoot it.

Let's look at an example using this package to log a domain event in what I would consider an acceptable fashion.

#### Example

Imagine we have a function like the following:

```go
type EmailSender interface {
	SendEmail(ctx context.Context, email, subject string) error
}

func SendWelcomeEmail(ctx context.Context, c EmailSender, email string) error {
	err := c.SendEmail(ctx, email, "Welcome to MyApp!")
	if err != nil {
		return fmt.Errorf("send welcome email failed: %s: %w", email, err)
	}
	return nil
}
```

The specifics of how this works aren't important, but what should be obvious is that there's no logging of any kind. My advice above was to log anything that has side-effects in your system. Sending an email is definitely a side-effect, so by extension we should have logging. The first pass I'm used to seeing (and historically, doing) might look something like this (note this snippet is using the `log` package from the standard library):

```diff
	if err != nil {
+		log.Printf("failed to send welcome email to %s because %v", email, err)
		return fmt.Errorf("send welcome email failed: %s: %w", email, err)
```

This nets you something that looks like:

```
2023/01/31 22:22:36 failed to send welcome email to me@example.com because gateway timed out
```

This isn't bad and it's certainly readable to me, but parsing this programmatically would be a nightmare. Let's look at the structured approach:

```diff
        if err != nil {
+		log.Error(ctx, "send email failed", log.F{
+			"email": email,
+			"error": err,
+			"type":  "welcome",
+		})
		return fmt.Errorf("send welcome email failed: %s: %w", email, err)
```

Now, your log looks like this (this would be a single line in practice, but for visibility I've pretty printed it here):

```json
{
  "level": "ERROR",
  "timestamp": "2023-01-31T22:25:58.141571-05:00",
  "message": "send email failed",
  "tags": {},
  "fields": {
    "email": "me@example.com",
    "error": "gateway timed out",
    "type": "welcome"
  }
}
```

Immediately we can see the parallels to the previous example, but now the same information is optimized for readability by machines. In fact, nearly every modern logging platform can ingest logs formatted in this way, going so far as to extract known fields into first-class concepts. One thing to note is the format of the `message` field's value. Rather than set this to something more specific like `"send welcome email failed"`, I would advise going with the approach here for the reason mentioned above: grepability. By opting for a less specific message and including qualifying information in fields, you create an implicit grouping of similar events. Say, for instance, your email service is unreachable. Rather than looking at error rates on a per template basis, you can simply look at the error rate for the implicit "send email failed" group. If you need metrics on just the welcome emails, you can just drill down into the provided fields.

Now that we've seen how this works in the error case, we've still got one thing to address: the error case isn't the only meaningful domain event. This is how I would log this function in a production app:

```go
func SendWelcomeEmail(ctx context.Context, c EmailSender, email string) error {
	err := c.SendEmail(ctx, email, "Welcome to MyApp!")
	if err != nil {
		log.Error(ctx, "send email failed", log.F{"email": email, "error": err, "type": "welcome"})
		return fmt.Errorf("send welcome email failed: %s: %w", email, err)
	}
	log.Info(ctx, "sent email", log.F{"email": email, "type": "welcome"})
	return nil
}
```

With a fully orchestrated application adding request-specific tags, this results in a helpful log that might look like the following:

```json
{"level":"ERROR","timestamp":"2023-01-31T22:47:20.790073-05:00","message":"send email failed","tags":{"env":"prod","service":"database","request_id":"ec6c7b62-565a-40d5-92f9-d87c9f22a5a2"},"fields":{"email":"me@example.com","error":"too many requests","type":"welcome"}}
{"level":"INFO","timestamp":"2023-01-31T22:47:25.793284-05:00","message":"sent email","tags":{"env":"prod","service":"database","request_id":"c3c4696f-3041-4800-9204-87e9f9c377f7"},"fields":{"email":"you@example.com","type":"welcome"}}
```

Beautiful, isn't it?

## References

- https://peter.bourgon.org/go-best-practices-2016/#logging-and-instrumentation
