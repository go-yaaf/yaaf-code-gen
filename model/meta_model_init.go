package model

// ini
func (m *MetaModel) initModel() {

}

// initialize model with basic types from yaaf-common
func (m *MetaModel) initModel2() {
	tf := NewClassInfo("TimeFrame")
	tf.PackageFullName = "base"
	tf.AddField("From", "Timestamp", "From Timestamp")
	tf.AddField("To", "Timestamp", "To Timestamp")
	m.AddClassInfo(tf)

	// Add time data point
	tdp := NewClassInfo("TimeDataPoint", "TimeDataPoint model represents a generic datapoint in time")
	tdp.PackageFullName = "base"
	tdp.IsGeneric = true
	tdp.AddField("Timestamp", "Timestamp", "Datapoint Timestamp")
	tdp.AddField("value", "T", "Generic value")
	//tdp.GenericTypes["T"] = "T"
	tdp.GenericTypes = append(tdp.GenericTypes, StringKeyValue{Key: "T", Value: "T"})

	m.AddClassInfo(tdp)

	// Add time series
	ts := NewClassInfo("TimeSeries", "TimeSeries is a set of data points over time")
	ts.PackageFullName = "base"
	ts.IsGeneric = true
	ts.AddField("Name", "string", "Name of the time series")
	ts.AddField("Range", "TimeFrame", "Range of the series (from ... to)")
	ts.AddField("Values", "TimeDataPoint<T>", "Series data points")
	ts.GetField("Values").IsArray = true
	//ts.GenericTypes["T"] = "T"
	ts.GenericTypes = append(ts.GenericTypes, StringKeyValue{Key: "T", Value: "T"})

	m.AddClassInfo(ts)

	be := NewClassInfo("BaseEntity", "Base class for all entities")
	be.PackageFullName = "base"
	be.IsVisible = false
	be.AddField("Id", "string", "Unique object Id")
	be.AddField("CreatedOn", "Timestamp", "When the object was created [Epoch milliseconds Timestamp]")
	be.AddField("UpdatedOn", "Timestamp", "When the object was last updated [Epoch milliseconds Timestamp]")
	m.AddClassInfo(be)

	bx := NewClassInfo("BaseEntityEx", "Extended Base class for all entities")
	bx.PackageFullName = "base"
	bx.IsVisible = false
	bx.AddField("Id", "string", "Unique object Id")
	bx.AddField("CreatedOn", "Timestamp", "When the object was created [Epoch milliseconds Timestamp]")
	bx.AddField("UpdatedOn", "Timestamp", "When the object was last updated [Epoch milliseconds Timestamp]")
	bx.AddField("Flag", "number", "Entity status flag")
	bx.AddField("Props", "Json", "List of custom properties")
	m.AddClassInfo(bx)

	// Add time data point
	tuple := NewClassInfo("Tuple", "Tuple is a generic key-value pair")
	tuple.PackageFullName = "base"
	tuple.IsGeneric = true
	tuple.AddField("Key", "K", "Generic key")
	tuple.AddField("Value", "T", "Generic value")

	tuple.GenericTypes = append(tuple.GenericTypes, StringKeyValue{Key: "K", Value: "K"})
	tuple.GenericTypes = append(tuple.GenericTypes, StringKeyValue{Key: "T", Value: "T"})

	//tuple.GenericTypes["K"] = "K"
	//tuple.GenericTypes["T"] = "T"
	m.AddClassInfo(tuple)

	cd := NewClassInfo("ColumnDef")
	cd.PackageFullName = "base"
	cd.AddField("Icon", "string", "Column field icon")
	cd.AddField("Name", "string", "Column field name")
	cd.AddField("Type", "string", "Column field type")
	cd.AddField("Format", "string", "Display format hint")
	cd.AddField("Sort", "number", "Sort order 0: no sort, 1: sort asc, -1 sort desc")
	cd.AddField("FilterOp", "number", "Filter operand")
	cd.AddField("Filter", "number | number[] | string | string[]", "Filter by value")
	m.AddClassInfo(cd)
}
