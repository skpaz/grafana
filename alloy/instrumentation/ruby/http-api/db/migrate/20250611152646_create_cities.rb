class CreateCities < ActiveRecord::Migration[8.0]
  def change
    create_table :cities do |t|
      t.string :name
      t.string :state
      t.string :country
      t.integer :founded
      t.integer :population
    end
  end
end
