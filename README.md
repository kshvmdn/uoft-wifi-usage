## uoft-wifi-usage

View wifi usage by building from around the University of Toronto St. George campus.

__Data source__: http://maps.wireless.utoronto.ca/stg/index.php

### Installation

  - Install with Go (assumes you have Go [installed](https://golang.org/doc/install) and [configured](https://golang.org/doc/install#testing)):

    ```
    $ go get github.com/kshvmdn/uoft-wifi-usage
    $ uoft-wifi-usage -help
    ```

  - Build from source (also requires you to have Go installed):

    ```
    $ git clone https://github.com/kshvmdn/uoft-wifi-usage
    $ cd uoft-wifi-usage
    $ make
    $ ./uoft-wifi-usage -help
    ```

### Usage

  - Run the program with `-help` for the help dialogue.
  
    ```
    $ uoft-wifi-usage -help
    Usage of uoft-wifi-usage:
    -buildings string
        Building IDs (leave empty to view all)
    -verbose
        Show detailed output
    ```

  - Example:

    ```
    $ uoft-wifi-usage | more
    Astronomy & Astrophysics Building, 10 connections
    665 Spadina Ave, 9 connections
    Lash Miller Chemistry Labs, 50 connections
    Bissell Bldg., 24 connections
    215 Huron St., 6 connections
    481 University Ave, 0 connections
    Rosebrugh Bldg., 31 connections
    121 St George St, 9 connections
    Northrop Frye Hall, 6 connections
    ...
    $ uoft-wifi-usage -buildings=0080
    Bahen Centre for Information Technology, 133 connections
    $ uoft-wifi-usage -buildings="0080,12345" -verbose
    2017/02/22 00:55:54 Building with ID `12345` doesn't exist.
    Bahen Centre for Information Technology - 134 active connections, 57 active access points (of 148)
    ```

  - Refer to [this](buildings.csv) for a list of building IDs.

### Contribute

This project is completely open source, feel free to open an issue or submit a pull request!

##### Improvements

- [ ] Sort output by building name
- [ ] Add option to filter by building name (currently only supports ID)
