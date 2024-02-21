// analysis functions

package draw_stats

import (
	"fmt"
	"time"

	"github.com/davecgh/go-spew/spew"
)

// analyse list of time events
func analyseTimeEvents(events []TimeEvent) TimeEventAnalysis {
	var totalTime time.Duration=time.Duration(0)

	for i := range events {
        totalTime=totalTime+events[i].Duration
	}

    return TimeEventAnalysis {
        TotalTime: totalTime,
        AverageTime: time.Duration(
            totalTime.Nanoseconds()/int64(len(events)),
        ),

        Events: events,
    }
}

// generate a tag breakdown from a list of events for a certain target tag.
// all events must have the tag. events without the tag will print out a warning and
// be excluded from the breakdown
func genTagBreakdown(events []TimeEvent,targetTag TagType) TagBreakdown {
    var keyedEvents TimeEventsByTagValue=groupEventsByTagValue(events,targetTag)

    var analysisDict TagValueAnalysisDict

    var tagValue TagValue
    for tagValue = range keyedEvents {
        analysisDict[tagValue]=analyseTimeEvents(keyedEvents[tagValue])
    }

    return TagBreakdown {
        Tag: targetTag,

        ValuesAnalysis: analysisDict,
        KeyedEvents: keyedEvents,
    }
}

// group events in a list of events by their value for a certain target tag
func groupEventsByTagValue(events []TimeEvent,targetTag TagType) TimeEventsByTagValue {
    // events keyed by tag value
    var keyedEvents TimeEventsByTagValue=make(TimeEventsByTagValue)

    for i := range events {
        var tagValue TagValue
        var exists bool

        tagValue,exists=events[i].Tags[targetTag]

        if !exists {
            fmt.Println("event that was missing target tag")
            fmt.Println("missing tag: ",targetTag)
            fmt.Println("event:")
            spew.Dump(events[i])
            continue
        }

        _,exists=keyedEvents[tagValue]

        if !exists {
            keyedEvents[tagValue]=[]TimeEvent{}
        }

        keyedEvents[tagValue]=append(keyedEvents[tagValue],events[i])
    }

    return keyedEvents
}