-- main = interact custom
main = interact processInput

custom :: String -> String
custom _ = show $ fuelRequired 1969

processInput input = show $ sum [perLine line | line <- lines input]

perLine :: String -> Integer
perLine line = fuel
  where fuel = fuelRequired val
        val = read line

fuelRequired :: Integer -> Integer
fuelRequired mass
  | simple == 0 = simple
  | otherwise = simple + fuelRequired simple
  where simple = fuelRequiredSimple mass

fuelRequiredSimple :: Integer -> Integer
fuelRequiredSimple mass
  | amt < 0 = 0
  | otherwise = amt
  where amt = mass `div` 3 - 2
