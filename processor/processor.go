package processor

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"

	"github.com/go-yaaf/yaaf-code-gen/model"
)

// Processor interface
type Processor interface {
	Start() error
}

// BaseProcessor parses proto files and generates abstract meta Model
type BaseProcessor struct {
	Output string
	Model  *model.MetaModel
}

// initialize model with basic types from yaaf-common
/*
func (p *BaseProcessor) initModel() {
	tf := model.NewClassInfo("TimeFrame")
	tf.PackageFullName = "base"
	tf.AddField("From", "Timestamp", "From Timestamp")
	tf.AddField("To", "Timestamp", "To Timestamp")
	p.Model.AddClassInfo(tf)

	// Add time data point
	tdp := model.NewClassInfo("TimeDataPoint", "TimeDataPoint model represents a generic datapoint in time")
	tdp.PackageFullName = "base"
	tdp.IsGeneric = true
	tdp.AddField("Timestamp", "Timestamp", "Datapoint Timestamp")
	tdp.AddField("value", "T", "Generic value")
	//tdp.GenericTypes["T"] = "T"
	tdp.GenericTypes = append(tdp.GenericTypes, model.StringKeyValue{Key: "T", Value: "T"})

	p.Model.AddClassInfo(tdp)

	// Add time series
	ts := model.NewClassInfo("TimeSeries", "TimeSeries is a set of data points over time")
	ts.PackageFullName = "base"
	ts.IsGeneric = true
	ts.AddField("Name", "string", "Name of the time series")
	ts.AddField("Range", "TimeFrame", "Range of the series (from ... to)")
	ts.AddField("Values", "TimeDataPoint<T>", "Series data points")
	ts.GetField("Values").IsArray = true
	//ts.GenericTypes["T"] = "T"
	ts.GenericTypes = append(ts.GenericTypes, model.StringKeyValue{Key: "T", Value: "T"})

	p.Model.AddClassInfo(ts)

	be := model.NewClassInfo("BaseEntity", "Base class for all entities")
	be.PackageFullName = "base"
	be.IsVisible = false
	be.AddField("Id", "string", "Unique object Id")
	be.AddField("CreatedOn", "Timestamp", "When the object was created [Epoch milliseconds Timestamp]")
	be.AddField("UpdatedOn", "Timestamp", "When the object was last updated [Epoch milliseconds Timestamp]")
	p.Model.AddClassInfo(be)

	bx := model.NewClassInfo("BaseEntityEx", "Extended Base class for all entities")
	bx.PackageFullName = "base"
	bx.IsVisible = false
	bx.AddField("Id", "string", "Unique object Id")
	bx.AddField("CreatedOn", "Timestamp", "When the object was created [Epoch milliseconds Timestamp]")
	bx.AddField("UpdatedOn", "Timestamp", "When the object was last updated [Epoch milliseconds Timestamp]")
	bx.AddField("Flag", "number", "Entity status flag")
	bx.AddField("Props", "Json", "List of custom properties")
	p.Model.AddClassInfo(bx)

	// Add time data point
	tuple := model.NewClassInfo("Tuple", "Tuple is a generic key-value pair")
	tuple.PackageFullName = "base"
	tuple.IsGeneric = true
	tuple.AddField("Key", "K", "Generic key")
	tuple.AddField("Value", "T", "Generic value")

	tuple.GenericTypes = append(tuple.GenericTypes, model.StringKeyValue{Key: "K", Value: "K"})
	tuple.GenericTypes = append(tuple.GenericTypes, model.StringKeyValue{Key: "T", Value: "T"})

	//tuple.GenericTypes["K"] = "K"
	//tuple.GenericTypes["T"] = "T"
	p.Model.AddClassInfo(tuple)

	cd := model.NewClassInfo("ColumnDef")
	cd.PackageFullName = "base"
	cd.AddField("Icon", "string", "Column field icon")
	cd.AddField("Name", "string", "Column field name")
	cd.AddField("Type", "string", "Column field type")
	cd.AddField("Format", "string", "Display format hint")
	cd.AddField("Sort", "number", "Sort order 0: no sort, 1: sort asc, -1 sort desc")
	cd.AddField("FilterOp", "number", "Filter operand")
	cd.AddField("Filter", "number | number[] | string | string[]", "Filter by value")
	p.Model.AddClassInfo(cd)
}
*/

// File copies a single file from src to dst
func (p *BaseProcessor) fileCopy(src, dst string) error {
	var err error
	var srcFd *os.File
	var dstFd *os.File
	var srcInfo os.FileInfo

	if srcFd, err = os.Open(src); err != nil {
		return err
	}
	defer func() {
		_ = srcFd.Close()
	}()

	if dstFd, err = os.Create(dst); err != nil {
		return err
	}
	defer func() {
		_ = dstFd.Close()
	}()

	if _, err = io.Copy(dstFd, srcFd); err != nil {
		return err
	}
	if srcInfo, err = os.Stat(src); err != nil {
		return err
	}
	return os.Chmod(dst, srcInfo.Mode())
}

// Dir copies a whole directory recursively
func (p *BaseProcessor) dirCopy(src string, dst string) error {
	var err error
	var fds []os.FileInfo
	var fileInfo os.FileInfo

	if fileInfo, err = os.Stat(src); err != nil {
		return err
	}

	if err = os.MkdirAll(dst, fileInfo.Mode()); err != nil {
		return err
	}

	if fds, err = ioutil.ReadDir(src); err != nil {
		return err
	}
	for _, fd := range fds {
		srcPath := path.Join(src, fd.Name())
		dstPath := path.Join(dst, fd.Name())

		if fd.IsDir() {
			if err = p.dirCopy(srcPath, dstPath); err != nil {
				fmt.Println(err)
			}
		} else {
			if err = p.fileCopy(srcPath, dstPath); err != nil {
				fmt.Println(err)
			}
		}
	}
	return nil
}
