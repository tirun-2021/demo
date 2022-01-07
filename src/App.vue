<template>
  <div>
    <div id="my_dataviz"></div>

    <!-- the data table -->
    <div id="table"></div>
  </div>
</template>

<script>
import * as d3 from "d3";
import axios from "axios";
export default {
  name: "App",
  methods: {
    draw() {
      var data = [];
      var columns = ["x1", "y1", "x2", "y2", "distance"];
      var table = d3.select("#table").append("table");
      var thead = table.append("thead");
      var tbody = table.append("tbody");

      // color,x1,y1,x2,y2
      // ["#eee",1121,3131,432432,43243]
      var xyStore = [];

      // append the header row
      thead
        .append("tr")
        .selectAll("th")
        .data(columns)
        .enter()
        .append("th")
        .text(function (column) {
          return column;
        });

      // set the dimensions and margins of the graph
      var margin = { top: 10, right: 30, bottom: 30, left: 60 },
        width = 460 - margin.left - margin.right,
        height = 400 - margin.top - margin.bottom;

      // append the svg object to the body of the page
      var svg = d3
        .select("#my_dataviz")
        .append("svg")
        .attr("width", width + margin.left + margin.right)
        .attr("height", height + margin.top + margin.bottom)
        .append("g")
        .attr("transform", "translate(" + margin.left + "," + margin.top + ")");

      // Add X axis --> it is a date format
      var x = d3.scaleLinear().domain([0, 100]).range([0, width]);
      svg
        .append("g")
        .attr("transform", "translate(0," + height + ")")
        .call(d3.axisBottom(x));

      // Add Y axis
      var y = d3.scaleLinear().domain([0, 100]).range([height, 0]);
      svg.append("g").call(d3.axisLeft(y));

      // This allows to find the closest X index of the mouse:
      var bisect = d3.bisector(function (d) {
        return d.x;
      }).left;

      // Create a rect on top of the svg area: this rectangle recovers mouse position
      var rect = svg
        .append("rect")
        .style("fill", "#eee")
        .style("pointer-events", "all")
        .attr("width", width)
        .attr("height", height);

      var g = svg.append("g");

      rect.on("click", function (event) {
        var xy = d3.pointer(event, g.node());
        if (xyStore.length == 0) {
          // no point exist
          var color = "#" + (((1 << 24) * Math.random()) | 0).toString(16);
          g.append("circle")
            .attr("r", 5)
            .attr("cx", xy[0])
            .attr("cy", xy[1])
            .attr("fill", color);
          xyStore.push(color);
          xyStore.push(xy[0]);
          xyStore.push(xy[1]);
        } else {
          // exist point
          // get exist color
          color = xyStore[0];
          g.append("circle")
            .attr("r", 5)
            .attr("cx", xy[0])
            .attr("cy", xy[1])
            .attr("fill", color);

          xyStore.push(xy[0]);
          xyStore.push(xy[1]);

          // add line
          g.append("line")
            .attr("x1", xyStore[1])
            .attr("y1", xyStore[2])
            .attr("x2", xyStore[3])
            .attr("y2", xyStore[4])
            .style("stroke", xyStore[0]);

          var x1_ = x.invert(xyStore[1]);
          var y1_ = y.invert(xyStore[2]);
          var x2_ = x.invert(xyStore[3]);
          var y2_ = y.invert(xyStore[4]);

          var FormData = require("form-data");
          var postData = new FormData();
          postData.append("x1", x1_);
          postData.append("y1", y1_);
          postData.append("x2", x2_);
          postData.append("y2", y2_);

          var config = {
            method: "post",
            url: "http://127.0.0.1:8080/ST_Distance",
            data: postData,
          };

          var result = 0;

          axios(config)
            .then(function (response) {
              console.log(JSON.stringify(response.data));
              var rsp = response.data;
              console.log(rsp);
              // clean
              data.push({
                x1: x1_,
                y1: y1_,
                x2: x2_,
                y2: y2_,
                distance: parseFloat(rsp.result),
              });
              console.log(data);
              xyStore = [];
              tabulate(data); // 2 column table
            })
            .catch(function (error) {
              console.log(error);
            });
        }
      });

      function tabulate(data) {
        // create a row for each object in the data
        var rows = tbody.selectAll("tr").data(data).enter().append("tr");

        // create a cell in each row for each column
        var cells = rows
          .selectAll("td")
          .data(function (row) {
            return columns.map(function (column) {
              return { column: column, value: row[column] };
            });
          })
          .enter()
          .append("td")
          .text(function (d) {
            return d.value;
          });

        return table;
      }
    },
  },
  mounted() {
    this.draw();
  },
};
</script>

<style>
#app {
  font-family: "Avenir", Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>
