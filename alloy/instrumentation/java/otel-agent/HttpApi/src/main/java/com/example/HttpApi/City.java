package com.example.HttpApi;

public class City {
  private String name;
  private String state;
  private String county;
  private Integer founded;
  private Integer population;

  public City(
    String name,
    String state,
    String county,
    Integer founded,
    Integer population
  ) {
    this.name = name;
    this.state = state;
    this.county = county;
    this.founded = founded;
    this.population = population;
  }

  public String getName() {
    return name;
  }

  public String getState() {
    return state;
  }

  public String getCounty() {
    return county;
  }

  public Integer getFounded() {
    return founded;
  }

  public Integer getPopulation() {
    return population;
  }

}