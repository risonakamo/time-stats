package main

import (
	"fmt"
	"path/filepath"
	"runtime"
	"time-stats/time_stats"
	datadir_v2 "time-stats/time_stats/data_dir2"

	"github.com/gofiber/fiber/v2"
)

func main() {
	var here string=getHereDir()

    var dataDir string=filepath.Join(here,"../data")
    var metadataFile string=filepath.Join(here,"../data/metadata_v2.yml")

    var app *fiber.App=fiber.New(fiber.Config{
        CaseSensitive:true,
        EnablePrintRoutes:false,
    })


    // --- apis ---
    // return list of available data files
    app.Get("/data-names",func (c *fiber.Ctx) error {
        var datalist []datadir_v2.DataFileInfo2=datadir_v2.ReadMetadataFileV2(metadataFile)

        return c.JSON(datalist)
    })

    // return data from a requested data file. generates analysis on the data file.
    app.Post("/get-data",func (c *fiber.Ctx) error {
        var body time_stats.GetDataRequest
        var e error=c.BodyParser(&body)

        if e!=nil {
            panic(e)
        }

        var fullFilePath string=filepath.Join(dataDir,body.Filename)

        var timeEvents []time_stats.TimeEvent
        timeEvents,e=time_stats.ParseSheetTsv(fullFilePath,true)

        if e!=nil {
            fmt.Println("failed to parse sheet tsv")
            return e
        }

        time_stats.AddDateTags(timeEvents)

        timeEvents=time_stats.FilterEvents(
            timeEvents,
            time_stats.TagFiltersListToDict(body.Filters),
        )

        if len(timeEvents)==0 {
            fmt.Println("filter resulted in no events")
            return fmt.Errorf("no events")
        }

        var analysis time_stats.TimeEventAnalysis=time_stats.AnalyseTimeEvents(timeEvents)

        var tagAnalysis time_stats.TagBreakdownsDict=time_stats.TagBreakdownForAllTags(timeEvents)

        var response time_stats.GetDataResponse=time_stats.GetDataResponse {
            TopAnalysis: analysis,
            TagsAnalysis: tagAnalysis,
        }

        return c.JSON(response)
    })

    // request to update a datafile. might do nothing if the datafile has no update url
    app.Post("/update-data",func(c *fiber.Ctx) error {
        var body time_stats.UpdateDataRequest
        var err error=c.BodyParser(&body)

        if err!=nil {
            panic(err)
        }

        var datafiles datadir_v2.MetadataYamlV2=datadir_v2.ReadMetadataFileV2(metadataFile)

        var datafile datadir_v2.DataFileInfo2
        datafile,err=datadir_v2.FindDataFile(body.Filename,datafiles)

        if err!=nil {
            panic(err)
        }

        err=datadir_v2.FetchDataFile(datafile,dataDir)

        if err!=nil {
            panic(err)
        }

        return nil
    })


    // --- static ---
    app.Static("/",filepath.Join(here,"../time-stats-web/build"))

    app.Listen(":4200")
}

// get directory of main function
func getHereDir() string {
    var selfFilepath string
    _,selfFilepath,_,_=runtime.Caller(0)

    return filepath.Dir(selfFilepath)
}