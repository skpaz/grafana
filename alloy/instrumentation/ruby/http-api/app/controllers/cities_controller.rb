class CitiesController < ApplicationController

  # display all cities
  def index
    @cities = City.all
    render json: @cities
  end

  # display specific city
  def show
    @city = City.find(params[:id])
    render json: @city
  end

  # add a new city
  def create
    @city = City.create(
      name: params[:name],
      state: params[:state],
      country: params[:country],
      founded: params[:founded],
      population: params[:population]
    )
    render json: @city
  end

end
