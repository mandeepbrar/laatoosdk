# Data Pipelines

This guide explains how to set up data pipelines in Laatoo for importing data from CSV files and other sources. Pipelines provide a flexible way to orchestrate data flows from source to destination with transformation steps in between.

## Overview

A data pipeline consists of three main components:

1. **Data Source (Importer)**: Reads raw data (CSV, database, API)
2. **Processors**: Chain of services that transform, validate, or map data
3. **Data Destination (Exporter)**: Saves processed data to storage

The pipeline also supports error handling and retry mechanisms for robust data processing.

## Configuration
 
 ### Pipeline Module
 
 Create a module configuration to define your pipeline.
 
 **File**: `config/server/modules/studentimport.yml`
 
 ```yaml
 plugin: datapipeline
 settings:
   datasource: studentcsv
   datadestination: studentsaver
   processors:
     - studentmapper
   pipelinename: studentimportpipeline
   errorprocessor: importerror
   retrysource: importerror
 ```
 
 **Settings:**
 - `datasource`: Input source service name
 - `datadestination`: Output destination service name
 - `processors`: List of transformation services (executed in order)
 - `pipelinename`: Unique pipeline identifier
 - `errorprocessor`: Service for handling errors
 - `retrysource`: Service for retrying failed entries
 
 ## Component Implementation
 
 Define pipeline components as services in your plugin (e.g., `dev/plugins/myplugin/config/server/services/`).
 
 ### 1. Data Source (CSV Importer)
 
 Reads data from a CSV file.
 
 **File**: `dev/plugins/studentmgmt/config/server/services/studentcsv.yml`
 
 ```yaml
 servicemethod: csvprocessor.CsvImporter
 importstorageservice: datadir
 csvinputfile: students.csv
 importcsvhasheaders: true
 ```
 
 **Configuration:**
 - `servicemethod`: Use built-in CSV importer
 - `importstorageservice`: Storage service where file is located
 - `csvinputfile`: CSV filename
 - `importcsvhasheaders`: `true` if CSV has header row
 
 ### 2. Processor (Field Mapper)
 
 Maps CSV columns to entity fields.
 
 **File**: `dev/plugins/studentmgmt/config/server/services/studentmapper.yml`
 
 ```yaml
 servicemethod: objectmapprocessor.MapToObjectProcessor
 object: studentmgmt.Student
 fieldmappings:
   StudentId: StudentID
   Name: FullName
   Email: EmailAddress
   DateOfBirth: DOB
   Department:
     field: DepartmentCode
     lookupfield: DepartmentCode
     dataservice: dataadapter.dataservice.departments.Department
 ```
 
 **Configuration:**
 - `object`: Target entity type
 - `fieldmappings`: Maps entity fields (keys) to CSV columns (values)
   - **Simple mapping**: `StudentId: StudentID` maps CSV column to entity field
   - **Lookup mapping**: `Department` performs lookup to find referenced entity
 
 ### 3. Data Destination (Exporter)
 
 Saves processed data to database.
 
 **File**: `dev/plugins/studentmgmt/config/server/services/studentsaver.yml`
 
 ```yaml
 servicemethod: datastoreprocessor.DataExporter
 exportdataservice: dataadapter.dataservice.studentmgmt.Student
 ```
 
 **Configuration:**
 - `servicemethod`: Use built-in data exporter
 - `exportdataservice`: Entity data service for saving
 
 ### 4. Error Handler
 
 Manages records that fail to process.
 
 **File**: `dev/plugins/studentmgmt/config/server/services/importerror.yml`
 
 ```yaml
 servicemethod: errorprocessors.MemoryErrorsProcessor
 ```
 
 Stores errors in memory for the duration of the process.

## Student Management Example
 
 ### Example 1: Import Students from CSV
 
 **CSV File** (`students.csv`):
 
 ```csv
 StudentID,FullName,EmailAddress,DOB,DepartmentCode
 S001,John Doe,john.doe@example.com,2000-01-15,CS
 S002,Jane Smith,jane.smith@example.com,1999-12-10,MATH
 S003,Bob Johnson,bob.j@example.com,2001-05-20,PHY
 ```
 
 **Pipeline Configuration** (`config/server/modules/studentimport.yml`):
 
 ```yaml
 plugin: datapipeline
 settings:
   datasource: studentcsv
   datadestination: studentsaver
   processors:
     - studentmapper
   pipelinename: studentimportpipeline
   errorprocessor: importerror
   retrysource: importerror
 ```
 
 **CSV Source** (`dev/plugins/studentmgmt/config/server/services/studentcsv.yml`):
 
 ```yaml
 servicemethod: csvprocessor.CsvImporter
 importstorageservice: datadir
 csvinputfile: students.csv
 importcsvhasheaders: true
 ```
 
 **Field Mapper** (`dev/plugins/studentmgmt/config/server/services/studentmapper.yml`):
 
 ```yaml
 servicemethod: objectmapprocessor.MapToObjectProcessor
 object: studentmgmt.Student
 fieldmappings:
   StudentId: StudentID
   Name: FullName
   Email: EmailAddress
   DateOfBirth: DOB
   Department:
     field: DepartmentCode
     lookupfield: DepartmentCode
     dataservice: dataadapter.dataservice.departments.Department
 ```
 
 **Department Lookup Explained:**
 - CSV contains `DepartmentCode` (e.g., "CS")
 - Mapper looks up Department entity by `DepartmentCode` field
 - Returns `storableref` to the Department entity
 - Student record gets proper Department reference
 
 **Data Saver** (`dev/plugins/studentmgmt/config/server/services/studentsaver.yml`):
 
 ```yaml
 servicemethod: datastoreprocessor.DataExporter
 exportdataservice: dataadapter.dataservice.studentmgmt.Student
 ```
 
 **Error Handler** (`dev/plugins/studentmgmt/config/server/services/importerror.yml`):
 
 ```yaml
 servicemethod: errorprocessors.MemoryErrorsProcessor
 ```
 
 ### Example 2: Import Course Enrollments
 
 **CSV File** (`enrollments.csv`):
 
 ```csv
 StudentID,CourseCode,Semester,EnrollmentDate
 S001,CS101,Fall2024,2024-08-15
 S001,MATH201,Fall2024,2024-08-15
 S002,CS101,Fall2024,2024-08-16
 ```
 
 **Enrollment Mapper** (`dev/plugins/studentmgmt/config/server/services/enrollmentmapper.yml`):
 
 ```yaml
 servicemethod: objectmapprocessor.MapToObjectProcessor
 object: studentmgmt.Enrollment
 fieldmappings:
   Semester: Semester
   EnrollmentDate: EnrollmentDate
   Student:
     field: StudentID
     lookupfield: StudentId
     dataservice: dataadapter.dataservice.studentmgmt.Student
   Course:
     field: CourseCode
     lookupfield: CourseCode
     dataservice: dataadapter.dataservice.courses.Course
 ```
 
 **Key features:**
 - Two lookups: Student by StudentID, Course by CourseCode
 - Both lookups return `storableref` objects
 - Enrollment entity gets proper references to Student and Course
 
 ### Example 3: Import Student Results
 
 **CSV File** (`results.csv`):
 
 ```csv
 StudentID,CourseCode,Semester,Grade,Score
 S001,CS101,Fall2024,A,95.5
 S001,MATH201,Fall2024,B+,87.0
 S002,CS101,Fall2024,A-,90.0
 ```
 
 **Results Mapper** (`dev/plugins/studentmgmt/config/server/services/resultsmapper.yml`):
 
 ```yaml
 servicemethod: objectmapprocessor.MapToObjectProcessor
 object: studentmgmt.StudentResult
 fieldmappings:
   Semester: Semester
   Grade: Grade
   Score: Score
   Student:
     field: StudentID
     lookupfield: StudentId
     dataservice: dataadapter.dataservice.studentmgmt.Student
   Course:
     field: CourseCode
     lookupfield: CourseCode
     dataservice: dataadapter.dataservice.courses.Course
 ```
 
 **Pipeline Configuration** (`config/server/modules/resultsimport.yml`):
 
 ```yaml
 plugin: datapipeline
 settings:
   datasource: resultscsv
   datadestination: resultssaver
   processors:
     - resultsmapper
     - resultsvalidator
   pipelinename: resultsimportpipeline
   errorprocessor: importerror
 ```
 
 **Custom Validator** (`dev/plugins/studentmgmt/config/server/services/resultsvalidator.yml`):
 
 ```yaml
 servicemethod: studentmgmt.ResultsValidator
 ```

You can create custom processors for validation logic:

```go
func (v *ResultsValidator) Invoke(ctx core.RequestContext) error {
    result := ctx.GetData().(*StudentResult)
    
    // Validate score range
    if result.Score < 0 || result.Score > 100 {
        return errors.ValidationError(ctx, "Score must be between 0 and 100")
    }
    
    // Validate grade matches score
    if !isValidGradeForScore(result.Grade, result.Score) {
        return errors.ValidationError(ctx, "Grade doesn't match score")
    }
    
    return nil
}
```

## Running the Pipeline

### Expose Pipeline via Channel

**File**: `config/server/channels/studentimport.yml`

```yaml
parent: root
method: POST
service: studentimportpipeline
path: /admin/import/students
staticvalues:
  retries: 1
permission: admin.import
```

**Channel properties:**
- `service`: Pipeline service name (from `pipelinename`)
- `path`: HTTP endpoint to trigger pipeline
- `staticvalues`: Static parameters like retry count
- `permission`: Required permission to run import

### Trigger the Pipeline

```bash
# Upload CSV file to storage first
curl -X POST http://localhost:6060/admin/import/students

# Or with retries
curl -X POST http://localhost:6060/admin/import/students \
  -H "Content-Type: application/json" \
  -d '{"retries": 3}'
```

## Available Components

### Importers

**CSV Importer:**
```yaml
servicemethod: csvprocessor.CsvImporter
importstorageservice: datadir
csvinputfile: data.csv
importcsvhasheaders: true
```

**Database Importer:**
```yaml
servicemethod: datastoreprocessor.DataImporter
importdataservice: dataadapter.dataservice.myentity
```

**REST API Importer:**
```yaml
servicemethod: restprocessor.RestImporter
importobject: mymodule.MyEntity
importrestendpoint: https://api.example.com/data
importlist: true
importmethod: GET
```

### Processors

**Object Mapper:**
```yaml
servicemethod: objectmapprocessor.MapToObjectProcessor
object: mymodule.MyEntity
fieldmappings:
  EntityField: CSVColumn
```

**Custom Processor:**
```yaml
servicemethod: mymodule.MyCustomProcessor
# Custom configuration options
```

### Exporters

**Database Exporter:**
```yaml
servicemethod: datastoreprocessor.DataExporter
exportdataservice: dataadapter.dataservice.myentity
```

**CSV Exporter:**
```yaml
servicemethod: csvprocessor.CsvExporter
exportstorageservice: datadir
csvoutputfile: output.csv
exportcsvhasheaders: true
```

**REST API Exporter:**
```yaml
servicemethod: restprocessor.RestExporter
exportrestendpoint: https://api.example.com/save
```

## Error Handling

Monitor and handle import errors:

```yaml
# Error processor configuration
servicemethod: errorprocessors.MemoryErrorsProcessor
```

Failed records are collected and can be:
- Logged for review
- Retried with corrected data
- Exported to error report

Access error information:

```go
errorService := ctx.GetService("importerror")
errors := errorService.GetErrors(ctx)

for _, err := range errors {
    log.Error(ctx, "Import failed", 
        "record", err.Record,
        "error", err.Error)
}
```

## Best Practices

### CSV File Preparation

1. **Use Headers**: Always include header row
2. **Consistent Encoding**: Use UTF-8 encoding
3. **Date Format**: Use ISO format (YYYY-MM-DD)
4. **Required Fields**: Ensure all required fields are present
5. **Data Validation**: Validate data before import

### Field Mapping

1. **Clear Column Names**: Use descriptive CSV column names
2. **Document Mappings**: Comment complex mappings
3. **Lookup Keys**: Ensure lookup fields are unique
4. **Default Values**: Provide defaults for optional fields
5. **Type Conversion**: Handle type mismatches gracefully

### Pipeline Design

1. **Modular Processors**: Keep processors focused and reusable
2. **Validation Early**: Validate data before expensive operations
3. **Error Recovery**: Implement retry logic for transient failures
4. **Logging**: Log pipeline execution for debugging
5. **Testing**: Test with sample data first

### Performance

1. **Batch Processing**: Import in batches for large files
2. **Index Lookups**: Ensure lookup fields are indexed
3. **Async Execution**: Run pipelines asynchronously
4. **Monitor Memory**: Watch memory usage for large imports
5. **Chunking**: Split very large files into chunks

## See Also

- [Entities Guide](entities.md) - Define entities for import
- [Services Guide](services.md) - Create custom processors
- [Workflows Guide](workflows.md) - Trigger pipelines from workflows
