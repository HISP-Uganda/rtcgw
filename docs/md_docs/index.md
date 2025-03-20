# Introduction
This documentation introduces the **RTC Gateway**, an application that was developed to support the integration of both eCHIS and LabXpert systems with DHIS2 (eCBSS). Rather than have the eCHIS and LabXpert systems technical teams implement the integration with DHIS2 separately, the RTC Gateway was developed to help obscure the DHIS2 internal logic of creating enrollments to DHIS 2 programs, adding events to a program stage and also updating the data values of the events.
This was achieved by creating simple JSON payloads that resonate well with what is collected from the eCHIS system and the results received from the LabXpert system.
# What is the RTC Gateway?
The RTC Gateway (rtcgw)  is a golang application that helps with the integration of:

- **eCHIS** with **DHIS2 (eCBSS)** by saving clients/ presumptive TB patients in the "Presumptive TB Program"
- **LabXpert** with **eCBSS** by saving results from the Genexpert machine for the presumptive patients in the "Presumptive TB Program"

The application is open source and its code can be found in the following repository:

[https://github.com/HISP-Uganda/rtcgw](https://github.com/HISP-Uganda/rtcgw/)

## How the RTCGW works?
### eCHIS - eCBSS Integration
1. The application works by receiving a payload from the eCHIS system. 
   - This payload has patient attributes and some details concerning their symptoms.
   - The format of the payload is already predetermined and must be adhered to. Refer to the integration guide for more on this.
   - Violations of the validation rules set for this payload result in an error
2. The above payload is unpacked and used to create an enrollment in the **TB Presumptive Program** if none exists. 
   - If an enrollment is present, the client's attributes and data values are simply updated.
3. Each payload has a special field called the **echis_patient_id** that is used to link a registered patient to the results coming from the LabXpert system.
   
### LabXpert - eCBSS Integration
1. The application receives a payload from the LabXpert system.
   - This payload has lab results for the client matching the **echis_patient_id** specified within.
   - The results are processed to match the values expected by the DHIS2 and are saved in.
     - The Presumptive TB Program
     - Positive results lead to an enrollment in the **Laboratory Program** and are saved as part of the event in its program stage
2. If the payload is sent again:
   - The application looks for an existing registration in the TB Presumptive Program to update the results. The same is done for the Lab Program for positive results

The processing of the payloads in either case has been made asynchronous in order not to overwhelm the application/service at peak times. In other words, the processing of these payloads is done in the background, but the requesting app gets an immediate response.

# Configuration

`rtcgw` uses a single yaml configuration file.

1. The default path is `/etc/rtcgw/rtcgw.yml`


You can use `rtcgw --help` to see a list of the configuration parameters and for more details on each option.

You will need to configure the following settings:

### Configuration file (Yaml)
> The configuration file is of a yaml format, with the following configurations:
Default path <kbd>/etc/rtcgw/rtcgw.yml</kbd>


| Tool/Component                      | Description                                                                  | Default                                                         |
|-------------------------------------|------------------------------------------------------------------------------|-----------------------------------------------------------------|
| **Databse Configurations**          |                                                                              |                                                                 |
| **uri**                             | URL describing how to connect to the rtcgw database                          | "**postgres://postgres:postgres@localhost/rtcgw?sslmode=disable**" |
| **Server Configurations**           |                                                                              |                                                                 |
| **host**                            | The hostname/ IP address for the host where the mfl-integrator is installed  | **localhost**                                                   |
| **http_port**                       | The port on which to run the mfl-integrator daemon                           | **9090**                                                        |
| **logdir**                          | The log directory for the application log files                              | **/var/log/rtcgw**                                              |
| **redis_address**                   | The Redid Address                                                            | **127.0.0.1:6379**                                              |
| **migrations_dir**                  | The migrations directory used to update DB schema                            | **/usr/share/rtcgw/db/migrations**                              |
| **templates_directory**             | The templates directory with documentation files                             | **/usr/share/rtcgw/docs/templates**                             |
| **static_directory**                | The Static directory                                                         | **/usr/share/rtcgw/docs/static**                                |
| **docs_directory**                  | The MD docs directory. Each md doc in this directory will be rendered        | **/usr/share/rtcgw/docs/md_docs**                               |
| **API Configurations**              |                                                                              |                                                                 |
| **dhis2_base_url**                  | The DHIS2 (ECBSS) base API URL                                               |                                                                 |
| **dhis2_user**                      | The DHIS2 API username                                                       |                                                                 |
| **dhis2_password**                  | The DHIS2 API user password                                                  |                                                                 |
| **dhis2_pat**                       | The DHIS2 API personal access token                                          |                                                                 |
| **dhis2_auth_method**               | The DHIS2 API authentication method                                          | **Basic**                                                       |
| **dhis2_tracker_program**           | The UID for the Presumptive TB Program                                       | **gjQIrstTQtl**                                                            |
| **dhis2_tracker_program_stage**     | The UID for the stage in the Presumptive TB Program                          | **ur2x5MjVfk7**                                                            |
| **dhis2_laboratory_program**        | The UID for the Laboratory Program in DHIS2                                  | **tu5n7t2P9QJ**                                                            |
| **dhis2_lab_program_stage**         | The UID for the stage in the Laboratory Program                              | **ghtfCYiCD4F**                                                            |
| **dhis2_search_attribute**          | The UID of the tracked entity attribute (the ECHISID) used to search clients | **fCctScv7UHr**                                                            |
| **dhis2_tracked_entity_types**      | The DHIS2 Tracked Entity Type representing patients in DHIS2                 | **aP2ziFSDvV4**                                                            |
| **API Configuration DHIS2 Mapping** |                                                                              |                                                                 |
| **Attributes**                      |                                                                              |                                                                 |
| **echis_patient_id**                | TE Attribute UID for the eCHIS patient ID                                    | **fCctScv7UHr**                                                 |
| **patient_gender**                  | TE Attribute UID for the Patient's Sex                                       | **GnL13HAVFOm**                                                 |
| **national_identification_number**  | TE Attribute UID for the Patient's NIN                                       | **IpM29RJ5pnG**                                                 |
| **patient_age_in_years**            | TE Attribute UID for the Patient's age in years                              | **Gy1jHsTp9P6**                                                 |
| **patient_age_in_months**           | TE Attribute UID for the Patient's age in months                             | **PeN6BxQaUkm**                                                 |
| **patient_age_in_days**             | TE Attribute UID for the Patient's age in days                               | **lEeXsdlXFxe**                                                 |
| **patient_phone**                   | TE Attribute UID for the Patient's Phone                                     | **kHjlSoKd1K1**                                                 |
| **patient_name**                    | TE Attribute UID for the Patient Name                                        | **jWjSY7cktaQ**                                                 |
| **client_category**                 | TE Attribute UID for Client Category                                         | **cZ0RMYYJWFO**                                                 |
| **Data Elements**                   |                                                                              |                                                                 |
| **cough**                           | The DataElement UID for Cough                                                | **phnhiuyDm3F**                                                 |
| **fever**                           | The DataElement UID for Fever                                                | **BtlM8ES6F3M**                                                 |
| **weight_loss**                     | The DataElement UID for Weight Loss                                          | **YDW9qk42Pvr**                                                 |
| **poor_weight_gain**                | The DataElement UID for Poor Weight Gain                                     |                                                                 |
| **excessive_night_sweat**           | The DataElement UID for Excessive Night Sweat                                | **iclIzowRYL1**                                                 |
| **results**                         | The DataElement UID for TB Results                                           | **pD0tc8UxyGg**                                                 |
| **results_date**                    | The DataElement UID for TB results Date                                      | **JCdqUjZZuvx**                                                 |
| **diagnosed**                       | The DataElement UID for TB Diagnosis                                         | **xgV6fQAETIf**                                                 |
| **lab_results**                                | The DataElement UID for TB Results in Lab Program                            | **G15dOa6fDYS**                                                 |
| **lab_results_date**                                | The DataElement UID for TB Results Date in Lab Program                       | **uCErZvwNSBL**                                                 |
| **lab_diagnosis**                                | The DataElement UID for TB Diagnosis in Lab Program                          | **lw0Hapx5FQ0**                                                 |
| **lab_sample_referred_from_community**                                | The DataElement UID for Sample Referred from Community - in lab program      | **Nf4Tz0J2vA6**                                                 |


### Sample Ymal Configuration file
```Yaml
database:
  uri: "postgres://postgres:postgres@localhost/rtcgw?sslmode=disable"

server:
  host: "localhost"
  http_port: 9292
  logdir: "/tmp"
  redis_address: "127.0.0.1:6379"
  migrations_dir: "file:///usr/share/rtcgw/db/migrations"
  templates_directory: "/usr/share/rtcgw/docs/templates"
  static_directory: "/usr/share/rtcgw/docs/static"
  docs_directory: "/usr/share/rtcgw/docs/md_docs"

api:
  dhis2_base_url: "https://tbl-ecbss-dev.health.go.ug/api/"
  dhis2_user: "admin"
  dhis2_password: "district"
  dhis2_pat: ""
  dhis2_auth_method: "Basic"
  dhis2_tracker_program: "gjQIrstTQtl"
  dhis2_tracker_program_stage: "ur2x5MjVfk7"
  dhis2_laboratory_program: "tu5n7t2P9QJ"
  dhis2_lab_program_stage: "ghtfCYiCD4F"
  dhis2_search_attribute: "fCctScv7UHr"
  dhis2_tracked_entity_type: "aP2ziFSDvV4"
  dhis2_mapping:
    attributes:
      echis_patient_id: "fCctScv7UHr"
      patient_gender: "GnL13HAVFOm"
      national_identification_number: "IpM29RJ5pnG"
      patient_age_in_years: "Gy1jHsTp9P6"
      patient_age_in_months: "PeN6BxQaUkm"
      patient_age_in_days: "lEeXsdlXFxe"
      patient_phone: "kHjlSoKd1K1"
      patient_name: "jWjSY7cktaQ"
      client_category: "cZ0RMYYJWFO"
    data_elements:
      cough: "phnhiuyDm3F"
      fever: "BtlM8ES6F3M"
      weight_loss: "YDW9qk42Pvr"
      poor_weight_gain: ""
      excessive_night_sweat: "iclIzowRYL1"
      results: "pD0tc8UxyGg"
      results_date: "JCdqUjZZuvx"
      diagnosed: "xgV6fQAETIf"
      lab_results: "G15dOa6fDYS"
      lab_results_date: "uCErZvwNSBL"
      lab_diagnosis: "lw0Hapx5FQ0"
      lab_sample_referred_from_community: "Nf4Tz0J2vA6"
```

The configuration file can also be hot reloaded, in other words, changes made to it do not require restarting the application.

The names of the **attributes** and **data_elements** are consistent with what is specified in the payloads received from eCHIS or LabXpert.

# Deployment
To simplify the deployment, binaries targeting different platforms (operating systems) have been released. In this section we will focus on deployment on a Linux distribution - specifically Ubuntu 22.04 LTS.

1. Release can be found [here.](https://github.com/HISP-Uganda/rtcgw/releases)
2. Each release tag will have a Linux, Windows and macOS binary
3. Additionally, we have provided a debian package to help with deployment on a Linux server

## Prerequisites
Before installing the debian package, the following prerequisite software is required:

1. PostgreSQL

## Installation

```bash
sudo dpkg -i rtcgw_1.0.0_amd64.deb
```
Once installed, there will be two systemd services created and enabled automatically:

1. `rtcgw`
    - This handles the HTTP requests from eCHIS and LabXpert
2. `rtcgw-workers`
    - This handles the background tasks

To start, stop, restart and view status of the application services use the following commands respectively:

The rtcgw service
```bash
sudo service rtcgw start

sudo service rtcgw stop

sudo service rtcgw restart

sudo service rtcgw status
```

The rtcgw-workers service
```bash
sudo service rtcgw-workers start

sudo service rtcgw-workers stop

sudo service rtcgw-workers restart

sudo service rtcgw-workers status
```

## Uninstallation

```bash
sudo dpkg --purge rtcgw
```

## Maintenance & troubleshooting tips
To troubleshoot any challenges with the application, it is recommended to view the logs.

The logs are typically located in the `/var/log/rtcgw` directory and are quite telling. Near all expected errors have been handled in the application and should be logged for troubleshooting purposes.