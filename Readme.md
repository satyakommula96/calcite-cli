<!--
{% comment %}
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to you under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
{% endcomment %}
-->

# Calcite CLI

Calcite CLI is a command-line interface for executing SQL queries using the Calcite server. It provides a prompt where you can enter queries and view the results.

## Build

```
go build
```

## Installation

To install Calcite CLI, you can use the following command:

```
go get github.com/satyakommula96/calcite-cli
```

## Usage

After installing Calcite CLI, you can run it using the following command:

calcite-cli [flags]


Flags:
- `--url`: Connection URL to the Calcite server (default: http://localhost:8080)
- `--serialization`: Serialization parameter (default: protobuf)
- `--enablePartitionPruning`: Enable Partition Pruning (default: true)
- `--distributedExecution`: Distributed Execution (default: false)

Once the Calcite CLI prompt starts, you can enter your SQL queries. To exit the prompt, type `exit`.

## Dependencies

This project uses the following third-party dependencies:
- [github.com/apache/calcite-avatica-go/v5](https://github.com/apache/calcite-avatica-go/v5)
- [github.com/olekukonko/tablewriter](https://github.com/olekukonko/tablewriter)
- [github.com/spf13/cobra](https://github.com/spf13/cobra)

That's the basic usage of the Calcite CLI. You can customize the connection URL and other parameters using command flags.

Happy querying!


## License

This project is licensed under the Apache License. See the [LICENSE](LICENSE) file for more information.
