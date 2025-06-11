Rails.application.routes.draw do
  resources :cities, only: [:index,:show,:create]
end
