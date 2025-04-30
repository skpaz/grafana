package com.example.HttpApi;

import java.util.Arrays;
import java.util.List;

import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RestController;

@RestController
public class HttpApiController {

  @GetMapping("/cities")
  public List<City> getCities() {
    return Arrays.asList(
      new City("Seattle","WA","King",1851,737015),
      new City("Portland","OR","Multnomah",1845,652503),
      new City("Los Angeles","CA","Los Angeles",1781,3898747),
      new City("Phoenix","AZ","Maricopa",1867,1608139)
    );
  }

  @PostMapping("/cities")
  public City createCity(@RequestBody City city) {
    return city;
  }

}
