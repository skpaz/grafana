# This file should ensure the existence of records required to run the application in every environment (production,
# development, test). The code here should be idempotent so that it can be executed at any point in every environment.
# The data can then be loaded with the bin/rails db:seed command (or created alongside the database with db:setup).
#
# Example:
#
#   ["Action", "Comedy", "Drama", "Horror"].each do |genre_name|
#     MovieGenre.find_or_create_by!(name: genre_name)
#   end

cities = City.create(
  [
    {
      "name": "Seattle",
      "state": "WA",
      "country": "King",
      "founded": 1851,
      "population": 737015
    },
    {
      "name": "Portland",
      "state": "OR",
      "country": "Multnomah",
      "founded": 1845,
      "population": 652503
    },
    {
      "name": "Los Angeles",
      "state": "CA",
      "country": "Los Angeles",
      "founded": 1781,
      "population": 3898747
    },
    {
      "name": "Phoenix",
      "state": "AZ",
      "country": "Maricopa",
      "founded": 1867,
      "population": 1608139
    }
  ]
)